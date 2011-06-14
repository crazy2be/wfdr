package json

import (
	"fmt"
	"json"
	"os"
)

func DecodeFromFile(filename string, object interface{}) (e os.Error) {
	file, e := os.Open(filename)
	if e != nil {
		fmt.Printf("Failed to open %s, %s.\n", filename, e)
		return
	}
	enc := json.NewDecoder(file)
	e = enc.Decode(object)
	if e != nil {
		fmt.Printf("Failed to decode JSON object: %s from file %s. (Object end up as %#v).\n", e, filename, object)
		return
	}
	return
}

func EncodeToFile(filename string, object interface{}) (e os.Error) {
	file, e := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0755)
	if e != nil {
		fmt.Printf("Failed to open %s, %s.\n", filename, e)
		return
	}
	enc := json.NewEncoder(file)
	e = enc.Encode(object)
	if e != nil {
		fmt.Printf("Failed to encode JSON object: %s. (Object represented as %#v).\n", e, object)
		return
	}
	return
}