package main

import (
	"fmt"
	"path"
	"strings"
)

type Path struct {
	root      string
	layout    string
	subfolder string
	extra     string
	file      string
}

// Calculates the path compontents and returns a new path, with the root set to the root, the subfolder set to the subfolder, and the other parts calculated automatically. Any missing sections become ".". The layout looks like this:
//     root/layout/subfolder/extra/file
// Note that the root, subfolder, and extra parts of the path could each contain multiple levels of neseting; there is no garenteered depth. Likewise, none of them are garenteered to have any depth whatsoever.
// Also worth noting is that each section is searched for sequentially, so /foo/blar with a root argument of /foo will result in blar being assumed as the layout.
func CalcPath(root, subfolder, fullpath string) (p *Path) {
	p = new(Path)
	// Clean all paths to avoid errors with ../fllo/blar causing problems... for the most part.
	root = path.Clean(root)
	fullpath = path.Clean(fullpath)
	subfolder = path.Clean(subfolder)
	
	p.root = root
	p.subfolder = subfolder
	
	// Nothing to calculate
	if len(fullpath) <= len(root) {
		return
	}
	secondPart := path.Clean(fullpath[len(root)+1:])
	if secondPart[0] == '/' {
		secondPart = secondPart[1:]
	}
	secondSlash := strings.Index(secondPart, "/")
	
	// No more /s, nothing to do.
	if secondSlash == -1 {
		p.layout = secondPart
		return
	}
	
	p.layout = secondPart[:secondSlash]
	
	thirdPart := secondPart[secondSlash+1:]
	
	fourthPart := thirdPart
	
	if IsSubPath(thirdPart, subfolder) {
		fourthPart = thirdPart[len(subfolder)+1:]
	}
	
	lastSlash := strings.LastIndex(fourthPart, "/")
	
	if lastSlash == -1 {
		p.extra = "."
							p.file = fourthPart
							return
	}
	
	p.file  = fourthPart[lastSlash+1:]
	p.extra = fourthPart[:lastSlash]
	
	return
}

// Figures out if /foo is a sub-path of /foo/blar. e.g.:
//    IsSubPath("/foo/blar", "/foo")
// returns true, and
//    IsSubPath("/foo", "/blar")
// returns false.
func IsSubPath(full, short string) bool {
	full = path.Clean(full)
	short = path.Clean(short)
	if len(full) < len(short) {
		// Can't possibly be a sub path.
		return false
	}
	fullToShort := full[:len(short)]
	//log.Println(fullToShort, short)
	if fullToShort == short {
		return true
	}
	return false
}

// Joins the path together and returns the result.
func (p *Path) Join() string {
	return p.JoinWithLayout(p.layout)
}

func (p *Path) JoinWithLayout(layout string) string {
	return p.JoinWithRootAndLayout(p.root, layout)
}

func (p *Path) JoinWithRoot(root string) string {
	return p.JoinWithRootAndLayout(root, p.layout)
}

func (p *Path) JoinWithRootAndLayout(root, layout string) string {
	return p.JoinWithRootLayoutAndSubfolder(root, layout, p.subfolder)
}

func (p *Path) JoinWithRootLayoutAndSubfolder(root, layout, subfolder string) string {
	joinedpath := path.Join(root, layout, subfolder, p.extra, p.file)
	//fmt.Println("Joined path:", path.Clean(joinedpath), "from parts", root, layout, subfolder, p.extra, p.file)
	return path.Clean(joinedpath)
}

// Kinda ugly old functions

func (p *Path) WithoutRoot() string {
	return path.Join(p.layout, p.Last())
}

func (p *Path) Last() string {
	return path.Join(p.extra, p.file)
}

func (p *Path) String() string {
	return fmt.Sprintf("%#v", p)
	//return path.Join(p.root, p.WithoutRoot())
}