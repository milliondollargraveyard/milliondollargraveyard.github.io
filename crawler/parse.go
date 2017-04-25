package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"sync/atomic"
)

var reTitle = regexp.MustCompile(`(?s)<title>(.*?)<`)
var reDomain = regexp.MustCompile(`>?.*domain.*<?`)

type htmlResult struct {
	Title          string
	MentionsDomain bool
}

func ParseHTML(r io.Reader) htmlResult {
	out := htmlResult{}

	s, err := ioutil.ReadAll(r)
	if err != nil {
		return out
	}

	matches := reTitle.FindSubmatch(s)
	if len(matches) < 2 {
		return out
	}
	out.Title = strings.TrimSpace(string(matches[1]))

	for _, matched := range reDomain.FindAll(s, -1) {
		if bytes.ContainsAny(matched, ";+") {
			continue
		}
		out.MentionsDomain = true
		break
	}

	return out
}

type CountingReader struct {
	io.Reader
	count uint64
}

func (r *CountingReader) Read(buf []byte) (int, error) {
	n, err := r.Reader.Read(buf)
	atomic.AddUint64(&r.count, uint64(n))
	return n, err
}

func (r *CountingReader) Count() uint64 {
	return atomic.LoadUint64(&r.count)
}
