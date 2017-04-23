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
		Status int    `json:"status,omitempty"`
		Error  string `json:"error,omitempty"`
		Size   int64  `json:"size,omitempty"`
		Title  string `json:"title,omitempty"`
	} `json:"response,omitempty"`
}

var skip = flag.Uint("skip", 0, "skip this many lines before processing")
var limit = flag.Int("limit", -1, "abort after processing this many lines")
var concurrency = flag.Uint("c", 1, "concurrency")

func main() {
	flag.Parse()

	// Read a JSON per line
	dec := json.NewDecoder(os.Stdin)

	client := &http.Client{Timeout: 10 * time.Second}

	if *concurrency < 1 {
		*concurrency = 1
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

	var i uint
	var s Site
	for {
		i++

		if i < *skip {
			continue
		}
		if *limit > 0 && uint(*limit) <= i-*skip {
			break
		}

		if err := dec.Decode(&s); err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("decode failed on line %d: %s", i, err)
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
				s.Response.Title = ParseTitle(resp.Body)
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
