package iomod

// Provides a few basic utility funcitons and types for man-in-the-middle byte counting.

import (
	"io"
	"os"
	"fmt"
)

type ReadCounter struct {
	Count int
	reader io.Reader
}

func (r *ReadCounter) Read(b []byte) (n int, e os.Error) {
	n, e = r.reader.Read(b)
	r.Count += n
	return
}

func NewReadCounter(reader io.Reader) (r* ReadCounter) {
	fmt.Println("What?")
	r = new(ReadCounter)
	r.reader = reader
	r.Count = 0
	return
}