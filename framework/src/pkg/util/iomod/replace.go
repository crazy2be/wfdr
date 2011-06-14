package iomod

import (
	"io"
	"os"
	//"fmt"
	"regexp"
)

type Replacer struct {
	From io.Reader
	Start *regexp.Regexp
	End []byte
}

// TODO: Make this more efficient. I should only have to copy the bytes once, if that.
func (r *Replacer) Read(p []byte) (n int, err os.Error) {
	//var length = len(p)
	// Allocate a buffer of cap(p) bytes and fill it from the input Reader.
	//var buf = make([]byte, length)
	n, err = r.From.Read(p)
	//fmt.Printf("%s", p)
	//fmt.Println(n, err)
	// Regexp replace!
	copy(p, r.Start.ReplaceAll(p, r.End))
	return
}

func NewReplacer(From io.Reader, start string, end string) (rep *Replacer) {
	var r, _ = regexp.Compile(start)
	rep = new(Replacer)
	rep.From = From
	rep.Start = r
	rep.End = []byte(end)
	return
	//return &Replacer{from, r, []byte(end)}
}
