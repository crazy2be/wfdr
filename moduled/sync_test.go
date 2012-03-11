package moduled

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var test1 []byte = []byte(`.foobar {
	color: red;
	font-family: sans-serif;
}
`)

var test1_mobile []byte = []byte(`.foobar {
	font-size: 10px;
}
`)

var test1_out_mobile []byte = append(test1, test1_mobile...)
var test1_out_desktop []byte = test1

func TestSync(t *testing.T) {
	source, err := ioutil.TempDir("", "moduled_source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(source)

	dest, err := ioutil.TempDir("", "moduled_dest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dest)

	cm, err := NewCacheMonitor(source, dest, "css")
	if err != nil {
		t.Fatal("Creating cache monitor:", err)
	}

	err = ioutil.WriteFile(path.Join(source, "test1.css"), test1, 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(path.Join(source, "test1_mobile.css"), test1_mobile, 0644)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Starting Sync")
	err = cm.Sync()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Checking mobile output")
	mobileb, err := ioutil.ReadFile(path.Join(dest, "mobile", "test1.css"))
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(mobileb, test1_out_mobile) {
		t.Errorf("Output file does not match expected. Got `%s`, expected `%s`.", mobileb, test1_out_mobile)
	}

	t.Log("Checking desktop output")
	desktopb, err := ioutil.ReadFile(path.Join(dest, "desktop", "test1.css"))
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(desktopb, test1_out_desktop) {
		t.Errorf("Output file does not match expected. Got `%s`, expected `%s`.", desktopb, test1_out_desktop)
	}
}

func TestDaemon(t *testing.T) {
	// TODO
}
