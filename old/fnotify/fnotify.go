package fnotify

// Provides a more convenient binding to the inotify API that utilizes callback functions rather than a polling loop, and is thus more condusive to a web-based enviroment, where most things are operating as callbacks anyway.

import (
	"exp/inotify"
	"fmt"
	"log"
	"sort"
	"strings"
	// Local imports
	//"wfdr/paths"
)

type modifiedHandler struct {
	path     string
	callback func(string)
}

type modifiedHandlerArray []modifiedHandler

func (this modifiedHandlerArray) Len() int {
	return len(this)
}
func (this modifiedHandlerArray) Less(i, j int) bool {
	return len(this[i].path) > len(this[j].path)
}
func (this modifiedHandlerArray) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

type Watcher struct {
	watcher          *inotify.Watcher
	modifiedHandlers modifiedHandlerArray
}

func NewWatcher() (*Watcher, error) {
	iwatcher, err := inotify.NewWatcher()
	// TODO: Polling fallback
	if err != nil {
		log.Fatal("Error initializing inotify:", err)
		return nil, err
	}
	mh := make(modifiedHandlerArray, 0)
	return &Watcher{iwatcher, mh}, nil
}

// The Listen-and-callback loop. Callbacks can still be registered after a call to this function. Usually called by your module in main() as
//     go watcher.Watch()
// A call to this function will block forever.
// TODO: Does this need to be single-threaded? E.g. should modified communicate with this loop?
func (w *Watcher) Watch() {
	for {
		select {
		case ev := <-w.watcher.Event:
			for _, handler := range w.modifiedHandlers {
				if strings.HasPrefix(ev.Name, handler.path) {
					fmt.Println(handler)
					handler.callback(ev.Name)
				}
			}
			log.Println("event:", ev)
			log.Println("handlers:", w.modifiedHandlers)
			//case addreq :=
		case err := <-w.watcher.Error:
			log.Println("error:", err)
		}
	}
}

// Register a callback for when files are modified.
func (w *Watcher) Modified(path string, callback func(string)) error {
	err := w.watcher.AddWatch(path, inotify.IN_MODIFY)
	if err != nil {
		log.Println("Error watching file", path, ":", err)
		return err
	}
	w.modifiedHandlers = append(w.modifiedHandlers, modifiedHandler{path, callback})
	sort.Sort(w.modifiedHandlers)
	return nil
}

//func Register(path string, flags uint32, callback func (*Event)) {

//}
