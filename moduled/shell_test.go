package moduled

import (
	"bytes"
	"fmt"
	"testing"
)

func eq(t *testing.T, args []string, exp []string) {
	if len(args) != len(exp) {
		t.Errorf("Length of args does not match expected length. Got %v (length %d), expected %v (length %d)", args, len(args), exp, len(exp))
	}
	match := true
	for i := range args {
		if args[i] != exp[i] {
			match = false
		}
	}
	if !match {
		t.Errorf("Arg sequence '%v' does not match expected arg sequence '%v'.", args, exp)
	}
}

func nerr(t *testing.T, cmd string, err error) {
	if err != nil {
		t.Errorf("Got unexpected error when inputing simple command string '%s': '%s'", cmd, err)
	}
}

type shellTest struct {
	in  string
	out []string
	err bool
}

var shellTests = []shellTest{
	{"command1 arg1 arg2\n", []string{"command1", "arg1", "arg2"}, false},
	// This test MUST come after the one above it, as it tests history autocompletion.
	{"\x1B[A\n", []string{"command1", "arg1", "arg2"}, false},
	{"\x7F\x7F\n", []string{""}, false},
	{"abc\x7F\x7F\n", []string{"a"}, false},
	{"abcdefh\x1B[C\x1B[Dg\n", []string{"abcdefgh"}, false},
	{"abc\x7F\x7F\x7F\x7F\n", []string{""}, false},
}

func TestShell(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	s := NewShell(in, out)

	for _, tst := range shellTests {
		fmt.Fprintf(in, tst.in)
		args, err := s.Prompt()
		if !tst.err {
			nerr(t, tst.in, err)
		}
		eq(t, args, tst.out)
	}
}
