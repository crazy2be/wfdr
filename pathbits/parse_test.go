package pathbits

import (
	"testing"
	"reflect"
)

type parsetest struct {
	in string
	ex Bits
}

var parsetests []parsetest = []parsetest{
	{"foo/bar/dude.css", Bits{"foo/bar/", "dude", nil, ".css"}},
	{"foo/bar/dude_linux_i386_mobile.css", Bits{"foo/bar/", "dude", []string{"linux", "i386", "mobile"}, ".css"}},
}

func TestParsepathbits(t *testing.T) {
	for i := range parsetests {
		pbt := parsetests[i]
		res := Parse(pbt.in)
		if !reflect.DeepEqual(res, &pbt.ex) {
			t.Errorf("When calling pathbits(%s), got %#v, expected %#v.", pbt.in, res, pbt.ex)
		}
	}
}