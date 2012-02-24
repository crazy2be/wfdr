package pathbits

import (
	"path"
	"strings"
)

type Bits struct {
	Path      string
	Name      string
	Modifiers []string
	Ext       string
}

func Parse(fullpath string) *Bits {
	pb := new(Bits)
	name := ""
	
	pb.Path, name = path.Split(fullpath)
	pb.Ext = path.Ext(name)
	name = name[:len(name)-len(pb.Ext)]
	
	pb.Modifiers = strings.Split(name, "_")
	pb.Name = pb.Modifiers[0]
	pb.Modifiers = pb.Modifiers[1:]
	if len(pb.Modifiers) == 0 {
		pb.Modifiers = nil
	}
	
	return pb
}

func Name(fullpath string) string {
	bits := Parse(fullpath)
	return bits.Name
}