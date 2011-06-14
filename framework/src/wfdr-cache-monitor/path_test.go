package main

import (
	"testing"
)

type CalcPathTest struct {
	root string
	subfolder string
	fullpath string
	expected *Path
}

var calcPathTests = []CalcPathTest{
	{
		"modules/events/css/",
		"",
		"modules/events/css/desktop/foo.css",
		&Path{
			"modules/events/css",
			"desktop",
			".", 
			".", 
			"foo.css",
		},
	},
	{
		"modules/events/js/",
		"",
		"modules/events/css/mobile/somefolder/foo.css",
		&Path {
			"modules/events/js",
			"mobile",
			".",
			"somefolder",
			"foo.css",
		},
	},
	{
		"cache/js/",
		"pages",
		"cache/js/desktop/pages/foo.css",
		&Path {
			"cache/js",
			"desktop",
			"pages",
			".",
			"foo.css",
		},
	},
}

func TestCalcPath(t *testing.T) {
	for _, test := range calcPathTests {
		result := CalcPath(test.root, test.subfolder, test.fullpath)
		// Ugly.
		fail := result.root != test.expected.root ||
						result.layout != test.expected.layout ||
						result.subfolder != test.expected.subfolder ||
						result.extra != test.expected.extra ||
						result.file != test.expected.file
		if fail {
			t.Errorf("Expected %#v\n Got %#v\n Test %#v", test.expected, result, test)
		}
	}
}

type IsSubPathTest struct {
	full string
	short string
	expected bool
}

var isSubPathTests = []IsSubPathTest {
	{ "/foo/blar", "/foo", true },
	{ "/foo", "/blar", false },
	{ "foo/blar", "foo/", true },
}

func TestIsSubPath(t *testing.T) {
	for _, test := range isSubPathTests {
		result := IsSubPath(test.full, test.short)
		if result != test.expected {
			t.Errorf("Expected %#v, got %#v on %#v.", test.expected, result, test)
		}
	}
}