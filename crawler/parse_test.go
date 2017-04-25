package main

import (
	"bytes"
	"testing"
)

func TestParseTitle(t *testing.T) {
	tests := []struct {
		Body  string
		Title string
	}{
		{
			Body: `<html>
				<head>
				<title>Broken HTML title
				<meta http-equiv="refresh" content="1; URL=foo">
				</head>
			`,
			Title: "Broken HTML title",
		},
		{
			Body: `<html>
				<head>
				<title>
				Sensible HTML title
				</title>
				</head>
				</html>
			`,
			Title: "Sensible HTML title",
		},
		{
			Body:  `<html><title>Easy title</title></html>`,
			Title: "Easy title",
		},
	}
	for i, testcase := range tests {
		html := bytes.NewBufferString(testcase.Body)
		if got, want := ParseHTML(html).Title, testcase.Title; got != want {
			t.Errorf("case %d: got %q; want %q", i, got, want)
		}
	}
}
