// Provides a utility functions for modification of byte steams
package iomod

import (
	"io"
	"os"
	"fmt"
)

// ReadCounter allows counting of the number of bytes read from a reader without control of the underlying reader.
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
	r = new(ReadCounter)
	r.reader = reader
	r.Count = 0
	return
}