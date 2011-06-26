package iomod

import (
	"io"
	"os"
	"regexp"
)

// Replacer provides a simple way to replace bytes in an incoming stream with any bytes you like. It uses a regexp internally, but it should be noted that regexps that cross Read() boundries will not match both sides of the boundry. For this reason, the longer the pattern, the less likely you are to get accurate results. Only works 100% reliably with single-byte replacements.
type Replacer struct {
	From io.Reader
	Start *regexp.Regexp
	End []byte
}

func (r *Replacer) Read(p []byte) (n int, err os.Error) {
	n, err = r.From.Read(p)
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
}
