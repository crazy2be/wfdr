// Simple line-based diff package, only works on plain text files. 
package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

type Diff struct {
	Line    int    // Line that was added or removed
	Added   bool   // Was the line added (true), or removed (false)?
	Content string // Line content after the change
}

func DiffReaders(original, end io.Reader) []Diff {
	origb, _ := ioutil.ReadAll(original)
	endb, _ := ioutil.ReadAll(end)

	DiffBytes(origb, endb)

	return nil
}

func DiffBytes(origb, endb []byte) []Diff {
	origl := bytes.Split(origb, []byte("\n"))
	endl := bytes.Split(endb, []byte("\n"))

	// Check if the streams are at all different
	// Do length check first for efficiency reasons.
	if len(origl) == len(endl) {
		if bytes.Equal(origb, endb) {
			// Bytes are equal!
			return nil
		}
	}

	for i, _ := range origl {
		if i >= len(endl) {
			fmt.Println("Out of range panic coming up!")
			fmt.Println(origl, endl, i)
		}
		if bytes.Equal(origl[i], endl[i]) {
			continue
		}
		// Search forward for the line
		for j := i; j < len(endl); j++ {
			if bytes.Equal(origl[i], endl[j]) {
				fmt.Println("Found match for line", i, "at line", j)
			}
		}
		for j := i; j >= 0; j-- {
			if bytes.Equal(origl[i], endl[j]) {
				fmt.Println("Found match for line", i, "at line", j)
			}
		}
	}
	return nil
}

func main() {
	f1 := "blar\nblah\nchicken"
	f2 := "blar\nchicken\nblah"
	DiffBytes([]byte(f1), []byte(f2))
}

// Longest common sequence implementation, see http://en.wikipedia.org/wiki/Longest_common_subsequence_problem
// WARNING: This uses an unoptimized algorthm, and will require massive amounts of memory if used on large byte arrays. It is reccommended that you do it by line.
//func lcsLength(
