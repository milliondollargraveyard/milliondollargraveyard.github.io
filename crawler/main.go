package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Site struct {
	Href     string `json:"href"`
	Coords   string `json:"coords"`
	Title    string `json:"title"`
	Response struct {
		Status     int    `json:"status,omitempty"`
		Error      string `json:"error,omitempty"`
		Size       int64  `json:"size,omitempty"`
		Title      string `json:"title,omitempty"`
		Redirected string `json:"redirected,omitempty"`
	} `json:"response,omitempty"`
}

var skip = flag.Int("skip", 0, "skip this many lines before processing")
var limit = flag.Int("limit", -1, "abort after processing this many lines")
var concurrency = flag.Int("c", 1, "concurrency")

const maxReadSize = 262144

func main() {
	flag.Parse()

	// Read a JSON per line
	dec := json.NewDecoder(os.Stdin)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	if *concurrency < 1 {
		*concurrency = 1
	}
	if *skip < 0 {
		*skip = 0
	}

	workers := make(chan struct{}, *concurrency)
	results := make(chan Site)
	wait := sync.WaitGroup{}

	go func() {
		// Writer
		enc := json.NewEncoder(os.Stdout)

		for s := range results {
			if err := enc.Encode(&s); err != nil {
				log.Fatalf("encode failed: %s", err)
			}
			wait.Done()
		}
	}()

	var i int
	var s Site
	for {
		i++

		if *limit > 0 && *limit <= i-*skip {
			log.Printf("stopping early due to limit: %d (limit=%d, skip=%d)", i, *limit, *skip)
			break
		}

		if err := dec.Decode(&s); err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("decode failed on line %d: %s", i, err)
		}

		if i < *skip {
			continue
		}

		workers <- struct{}{}
		go func(s Site) {
			wait.Add(1)
			log.Printf("[%d] Processing %s", i, s.Href)

			resp, err := client.Get(s.Href)
			if err != nil {
				s.Response.Error = err.Error()
			} else {
				s.Response.Status = resp.StatusCode
				if resp.ContentLength >= 0 {
					s.Response.Size = resp.ContentLength
				}
				url := resp.Request.URL.String()
				if url != s.Href {
					s.Response.Redirected = url
				}
				readcount := &ReaderCounter{Reader: resp.Body}
				s.Response.Title = ParseTitle(readcount)
				if s.Response.Redirected == "" && s.Response.Size == 0 {
					// Consume the body until completion to measure the body
					for {
						if _, err := readcount.Read(nil); err == io.EOF {
							break
						}
						if readcount.Count() > maxReadSize {
							break
						}
					}
					s.Response.Size = int64(readcount.Count())
				}
				resp.Body.Close()
			}

			results <- s
			<-workers
		}(s)
	}

	log.Printf("shutting down after line %d", i)
	wait.Wait()
	close(results)
}
