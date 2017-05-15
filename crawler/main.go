package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Archive struct {
	*ArchiveResult
	Error string `json:"error,omitempty"`
}

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
		Squatter   bool   `json:"squatter,omitempty"`
	} `json:"response,omitempty"`
	Archive *Archive `json:"archive,omitempty"`
}

// Records returns a CSV-compatible flat output
func (s *Site) Records() []string {
	return []string{
		s.Href, s.Coords, s.Title,
		fmt.Sprintf("%d", s.Response.Status),
		s.Response.Error,
		fmt.Sprintf("%d", s.Response.Size),
		s.Response.Title,
		s.Response.Redirected,
		fmt.Sprintf("%t", s.Response.Squatter),
	}
}

var skip = flag.Int("skip", 0, "skip this many lines before processing")
var limit = flag.Int("limit", -1, "abort after processing this many lines")
var concurrency = flag.Int("c", 1, "concurrency")
var exportCSV = flag.Bool("csv", false, "export in csv format")
var loadArchive = flag.Bool("archive", false, "load snapshot from archive.org")

const maxReadSize = 262144

func main() {
	flag.Parse()

	// Read a JSON per line
	var dec = json.NewDecoder(os.Stdin)
	var enc Encoder

	if *exportCSV {
		enc = &CSVEncoder{csv.NewWriter(os.Stdout)}
	} else {
		enc = json.NewEncoder(os.Stdout)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	if *concurrency < 1 {
		*concurrency = 1
	}
	if *skip < 0 {
		*skip = 0
	}

	archiveMutex := sync.Mutex{}
	workers := make(chan worker, *concurrency)
	for i := 0; i < *concurrency; i++ {
		w := worker{
			client: client,
		}
		if *loadArchive {
			w.loadArchive = true
			w.lockArchive = &archiveMutex
		}
		workers <- w
	}
	results := make(chan Site)
	wait := sync.WaitGroup{}

	go func() {
		// Writer
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

		w := <-workers
		go func(s Site) {
			wait.Add(1)
			log.Printf("[%d] Processing %s", i, s.Href)

			results <- w.Process(s)
			workers <- w
		}(s)

		i++
	}

	log.Printf("shutting down after line %d", i)
	wait.Wait()
	close(results)
	close(workers)
}

type worker struct {
	client      *http.Client
	loadArchive bool
	lockArchive sync.Locker
}

// Process takes a Site query and returns an augmented Site
func (w *worker) Process(s Site) Site {
	if !w.loadArchive {
		w.lockArchive.Lock()

		archive := ArchiveClient{w.client}
		// The auction ended in Jan 2006, so we fetch archives after Feb
		results, err := archive.History(s.Href, Param{"limit", "1"}, Param{"from", "200602"})
		arch := &Archive{}
		if err != nil {
			arch.Error = err.Error()
			log.Printf("archive error: %s", err)
		} else if len(results) > 0 {
			arch.ArchiveResult = &results[0]
		} else {
			arch = nil
		}
		s.Archive = arch
		w.lockArchive.Unlock()
	}

	resp, err := w.client.Get(s.Href)
	if err != nil {
		s.Response.Error = err.Error()
		return s
	}
	defer resp.Body.Close()

	// Fill in the response
	s.Response.Status = resp.StatusCode
	if resp.ContentLength >= 0 {
		s.Response.Size = resp.ContentLength
	}

	// Is it a redirect?
	url := resp.Request.URL.String()
	if url != s.Href {
		s.Response.Redirected = url
	}

	// What's in the body?
	readcount := &CountingReader{Reader: resp.Body}
	parsed := ParseHTML(readcount)
	s.Response.Title = parsed.Title
	s.Response.Squatter = parsed.MentionsDomain

	// Should we try to guess the content size?
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

	return s
}
