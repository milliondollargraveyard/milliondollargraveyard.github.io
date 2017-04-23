package main

import (
	"io"
	"strings"
	"sync/atomic"

	"golang.org/x/net/html"
)

func ParseTitle(r io.Reader) (title string) {
	doc, err := html.Parse(r)
	if err != nil {
		return
	}

	var traverse func(*html.Node) bool
	traverse = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = n.FirstChild.Data
			}
			return false
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if !traverse(c) {
				return false
			}
		}
		return true
	}
	traverse(doc)

	if title == "" {
		return
	}

	parts := strings.SplitN(title, "<", 2)
	if len(parts) == 2 {
		title = parts[0]
	}

	return strings.TrimSpace(title)
}

type ReaderCounter struct {
	io.Reader
	count uint64
}

func (r *ReaderCounter) Read(buf []byte) (int, error) {
	n, err := r.Reader.Read(buf)
	atomic.AddUint64(&r.count, uint64(n))
	return n, err
}

func (r *ReaderCounter) Count() uint64 {
	return atomic.LoadUint64(&r.count)
}
