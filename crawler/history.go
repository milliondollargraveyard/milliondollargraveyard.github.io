package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Param struct {
	Key   string
	Value string
}

var defaultUserAgent = "milliondollarcrawler/0.1"

const archiveTimestampLayout = "20060102150405"
const archiveAPI = "https://web.archive.org/cdx/search/cdx"

type ArchiveResult struct {
	URL       string    `json:"url"`
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	Length    int       `json:"length"`
}

func (h *ArchiveResult) Request(method string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, h.URL, body)
}

type ArchiveClient struct {
	*http.Client
}

func (f *ArchiveClient) readURL(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", defaultUserAgent)

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

// History returns a slice of ArchiveResults based on Params like `limit`, `from`, `to`, etc.
// See https://github.com/internetarchive/wayback/blob/master/wayback-cdx-server/README.md
func (f *ArchiveClient) History(URL string, params ...Param) ([]ArchiveResult, error) {
	p := url.Values{}
	p.Add("url", URL)
	p.Add("output", "json")
	for _, param := range params {
		p.Add(param.Key, param.Value)
	}

	body, err := f.readURL(archiveAPI + "?" + p.Encode())
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var csv [][]string
	dec := json.NewDecoder(body)
	if err := dec.Decode(&csv); err != nil {
		return nil, err
	}

	r := []ArchiveResult{}

	if len(csv) < 2 {
		return r, nil
	}

	// Skip the first row, the headers
	// ["urlkey","timestamp","original","mimetype","statuscode","digest","length"]
	for _, vals := range csv[1:] {
		timestampRaw, statusStr, lengthStr := vals[1], vals[4], vals[6]
		if statusStr != "200" {
			continue
		}

		ts, err := time.Parse(archiveTimestampLayout, timestampRaw)
		if err != nil {
			continue
		}

		status, _ := strconv.Atoi(statusStr)
		length, _ := strconv.Atoi(lengthStr)

		r = append(r, ArchiveResult{
			Timestamp: ts,
			URL:       "https://web.archive.org/web/" + timestampRaw + "/" + URL,
			Status:    status,
			Length:    length,
		})
	}

	return r, nil
}
