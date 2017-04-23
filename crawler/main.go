package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
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

var skip = flag.Int("skip", 0, "skip this many lines before processing")
var limit = flag.Int("limit", -1, "abort after processing this many lines")

func main() {
	flag.Parse()

	// Read a JSON per line
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)

	client := &http.Client{Timeout: 10 * time.Second}

	var i int
	var s Site
	for {
		i++

		if i < *skip {
			continue
		}
		if *limit > 0 && *limit <= i-*skip {
			break
		}

		if err := dec.Decode(&s); err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("decode failed on line %d: %s", i, err)
		}

		log.Printf("[%d] Processing %s", i, s.Href)

		resp, err := client.Get(s.Href)
		s.Response.Status = resp.StatusCode
		if resp.ContentLength >= 0 {
			s.Response.Size = resp.ContentLength
		}
		if err != nil {
			s.Response.Error = err.Error()
		} else {
			s.Response.Title = ParseTitle(resp.Body)
			resp.Body.Close()
		}

		if err := enc.Encode(&s); err != nil {
			log.Fatalf("encode failed on line %d: %s", i, err)
		}
	}

	log.Printf("done after line %d", i)
}
