package util

// This file provides custom CSV readers and writers, using the stdlib
// CSV reader and writer.
// Since most of the tools read a format (int, int[, something, ...])
// it is useful to have the int to string conversons in one place so
// I don't repeat myself too much.

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Write is a CSV custom writer type.
type Writer struct {
	w *csv.Writer
}

// NewWriter creates a new writer.
func NewWriter(f io.Writer) *Writer {
	return &Writer{csv.NewWriter(f)}
}

func (w *Writer) Flush() error {
	w.w.Flush()
	return w.w.Error()
}

// Write writes a CSV record a,b,rubbish...
func (w *Writer) Write(a, b int, rubbish []string) error {
	line := make([]string, len(rubbish)+2)
	line[0] = strconv.Itoa(a)
	line[1] = strconv.Itoa(b)
	copy(line[2:], rubbish)
	return w.w.Write(line)
}

// Reader is a custom CSV Reader.
type Reader struct {
	r *csv.Reader
}

// NewReader creates new CSV reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{csv.NewReader(r)}
}

// Read reads an element.
func (r *Reader) Read() (int, int, []string, error) {
	record, err := r.r.Read()
	if err != nil {
		return 0, 0, nil, err
	}
	if len(record) < 2 {
		return 0, 0, record, fmt.Errorf("Not enough data in the record.")
	}
	record[0] = strings.TrimSpace(record[0])
	record[1] = strings.TrimSpace(record[1])
	a, err := strconv.Atoi(record[0])
	if err != nil {
		return 0, 0, record, err
	}
	b, err := strconv.Atoi(record[1])
	if err != nil {
		return 0, 0, record, err
	}
	if len(record) == 2 {
		return a, b, nil, nil
	}
	return a, b, record[2:], nil
}
