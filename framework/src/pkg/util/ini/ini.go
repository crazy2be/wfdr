// A simple package to read ini files (actually just name=value files, [sections] are not supported yet because we have no need for them.
package ini

import (
	"io/ioutil"
	"bytes"
	"path"
	"fmt"
	"os"
	// Local imports
	"util/dlog"
)

// Load an ini file. Pass a filename, returns a map of all of the name=value pairs within the file, and an error if applicable. 
func Load(filename string) (settings map[string]string, e os.Error) {
	settings = make(map[string]string, 10)
	fileContents, e := ioutil.ReadFile(filename)
	if e != nil {
		dlog.Printf("Failed to read %s, %s.\n", filename, e)
		return
	}
	lines := bytes.Split(fileContents, []byte("\n"), -1)
	for _, line := range(lines) {
		splitLine := bytes.Split(line, []byte("="), 2)
		if len(splitLine) < 2 {
			//No = in line or value is blank. Skipping.
			break;
		}
		//fmt.Printf("Setting, value: %s, %s\n", splitLine[0], splitLine[1])
		settings[string(splitLine[0])] = string(splitLine[1])
	}
	return
}

// Save a map of settings to filename in name=value format. Returns nil on success, error otherwise.
func Save(filename string, settings map[string]string) (e os.Error) {
	dirname, _ := path.Split(filename)
	err := os.MkdirAll(dirname, 0755)
	if err != nil {
		dlog.Println("Error creating directory to save ini file", filename, ":", err)
		return err
	}
	
	file, e := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if e != nil {
		dlog.Printf("Failed to open %s, %s.\n", filename, e)
		return
	}
	for key, value := range(settings) {
		fmt.Fprintf(file, "%s=%s\n", key, value)
	}
	return
}