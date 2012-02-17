package paths

// Utility library that provides various functions for matching and sorting lists of paths, in order to use http-style pattern matching (see the Match() function).
// This package is useless until generics are implemented...

import (
	"fmt"
	"sort"
	//"path"
	//"reflect"
	//"log"
)

// Showing it not working...
// type blar struct {
// 	pathp1 string
// 	pathp2 string
// }
// 
// func (b *blar) Path() string {
// 	return path.Join(b.pathp1, b.pathp2)
// }
// 
// func main() {
// 	blars := make([]blar, 10)
// 	Sort(blars)
// }

// A Pather is a type that implements Path(), such that this package can use it.
type Pather interface {
	Path() string
}

type Pathers []Pather

func (this Pathers) Len() int {
	return len(this)
}
func (this Pathers) Less(i, j int) bool {
	return len(this[i].Path()) > len(this[j].Path())
}
func (this Pathers) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// Sorts the paths accoding to length, such that more specific patterns get considered first.
func Sort(paths []Pather) {
	pathers := Pathers(paths)
	sort.Sort(pathers)
	return
}

// Uses HTTP-style path matching. /foo/ will match /foo and /foo/blar, while /foo will only match /foo. Returns true if path matches pattern.
func Match(pattern, path string) bool {
	if len(pattern) == 0 {
		// should not happen
		return false
	}
	n := len(pattern)
	if pattern[n-1] != '/' {
		return pattern == path
	}
	return len(path) >= n && path[0:n] == pattern
}

// Searches through the list of patterns (should be sorted via a call to Sort() first), and returns the first one that matches. If no match is located, nil is returned.
func FindMatch(patterns []Pather, path string) Pather {
	for _, pattern := range patterns {
		if Match(pattern.Path(), path) {
			// TODO: Caching
			fmt.Println("Found match for path", path, ":", pattern.Path())
			return pattern
		}
	}
	return nil
}
