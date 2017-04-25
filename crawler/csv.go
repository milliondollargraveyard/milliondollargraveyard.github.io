package main

import (
	"encoding/csv"
	"errors"
)

type Encoder interface {
	Encode(v interface{}) error
}

type CSVEncodable interface {
	Records() []string
}

type CSVEncoder struct {
	*csv.Writer
}

func (w *CSVEncoder) Encode(v interface{}) error {
	row, ok := v.(CSVEncodable)
	if !ok {
		return errors.New("not encodeable")
	}
	err := w.Writer.Write(row.Records())
	w.Writer.Flush()
	return err
}
