package main

import (
	"io"

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
			title = n.FirstChild.Data
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

	return
}
