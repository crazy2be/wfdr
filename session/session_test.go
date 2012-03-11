package session

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGenID(t *testing.T) {
	id := GenID(nil)
	if id == "" {
		t.Fatal("GenID() returned empty string!")
	}
	id2 := GenID([]byte("dklfjalkdsjfalk;jdsfl"))
	for i := 0; i < 100; i++ {
		if id == id2 {
			t.Fatal("GenID() returned two duplicate ids!")
		}
		id = id2
		id2 = GenID(nil)
	}
	t.Logf("id1: %s\nid2: %s\n", id, id2)
}

func TestFSManager(t *testing.T) {
	dir, err := ioutil.TempDir("", "session")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	fsm, err := NewFSManager(dir)
	if err != nil {
		t.Fatal("Unexpected error creating FSManager:", err)
	}

	id := GenID(nil)
	val, err := fsm.Get(id, "foobar")
	if err == nil {
		t.Error("Expected error when attempting to get a value for a non-existant ID")
	}
	if val != "" {
		t.Error("Expected empty string when attempting to get a value for a non-existant id and value.")
	}

	err = fsm.Set(id, "foobar", "testvalue")
	if err != nil {
		t.Error("Unexpected error when attempting to set a value 'foobar' for a newly created id.")
	}

	val, err = fsm.Get(id, "foobar")
	if err != nil {
		t.Error("Unexpected error when attempting to fetch a previously set key 'foobar'.")
	}
	if val != "testvalue" {
		t.Errorf("Set 'foobar' to 'testvalue', got back unexpected value of '%s'.", val)
	}

}
