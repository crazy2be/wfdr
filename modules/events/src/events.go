package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
	// Local imports
	"util/pages"
	tmpl "util/template"
	//"util/perms"
)

type Time struct {
	time.Time
}

func (t *Time) String() string {
	if t == nil {
		return ""
	}
	return t.Format(timeFormat)
}

func (t *Time) SimpleString() string {
	weekday := "Unknown"
	switch t.Weekday() {
	case 0:
		weekday = "Sunday"
	case 1:
		weekday = "Monday"
	case 2:
		weekday = "Tuesday"
	case 3:
		weekday = "Wendesday"
	case 4:
		weekday = "Thursday"
	case 5:
		weekday = "Friday"
	case 6:
		weekday = "Saturday"
	}
	suffix := "th"
	switch t.Day() % 10 {
	case 1:
		suffix = "st"
	case 2:
		suffix = "nd"
	case 3:
		suffix = "rd"
	}
	day := strconv.Itoa(t.Day())
	return weekday + " the " + day + suffix
}

type Event struct {
	ID         string // Unique
	Title      string
	Desc       string // Should be a []byte, but mustache doesn't seem to display that how i'd like..
	Link       string
	Time       time.Time
	Img        string
	Importance int
	PageData   *pages.PageData // nil on the main page.
}

// Allows events to be sorted by importance, removing wierd floating issues that exist otherwise.
type EventArray []*Event

func (e EventArray) Len() int {
	return len(e)
}
func (e EventArray) Less(i, j int) bool {
	return e[i].Importance < e[j].Importance
}
func (e EventArray) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type Filler struct {
	Importance int
}

type EventsObj struct {
	Events []*Event
	Filler *Filler
}

func getEvents() (events map[string]*Event) {
	// TODO: Caching
	events = make(map[string]*Event, 10)
	// Open the directory
	eventDir, _ := os.Open("data/events/event")
	// WARNING: Reads ALL events from disk, could be an issue if we have loads of events...
	eventFileNames, e := eventDir.Readdirnames(-1)
	if e != nil {
		fmt.Println(e)
	}
	for _, eventFileName := range eventFileNames {
		event := new(Event)
		event.ID = eventFileName[:len(eventFileName)-len(".json")]

		e := event.Load()
		if e != nil {
			continue
		}

		events[event.ID] = event
	}
	return
}

func saveEvents(events map[string]*Event) (e error) {
	for _, event := range events {
		e = event.Save()
		if e != nil {
			return
		}
	}
	return
}

// Loads the event data from disk. Caller MUST set event.ID before calling.
func (event *Event) Load() (e error) {
	fmt.Println(event.ID)
	eventFile, e := os.Open("data/events/event/" + event.ID + ".json")
	if e != nil {
		fmt.Println("Error opening event JSON:", e)
		return
	}

	decoder := json.NewDecoder(eventFile)
	e = decoder.Decode(&event)
	if e != nil {
		fmt.Println("Error decoding event JSON:", e)
		return
	}
	fmt.Println(event)
	return
}

// Loads the page data and content. Seperate from Load() because it requires more time than just loading the basic data, and is not required on the main page.
func (event *Event) LoadPage(r *http.Request) (e error) {
	pageData, e := pages.GetPageData("events/"+event.ID, r)
	event.PageData = pageData
	return
}

// Saves the event data to disk. Caller must set event.ID before calling.
func (event *Event) Save() (e error) {
	id := event.ID
	eventFile, e := os.Create("data/events/event/" + id + ".json")
	if e != nil {
		fmt.Println("Error saving events:", e)
		return
	}
	encoder := json.NewEncoder(eventFile)
	e = encoder.Encode(event)
	if e != nil {
		fmt.Println("Error encoding event JSON:", e)
		return
	}
	return
}

func (event *Event) Delete() (e error) {
	if len(event.ID) <= 0 {
		return errors.New("Invalid ID!")
	}
	e = os.Remove("data/events/event/" + event.ID + ".json")
	return
}

// Makes a filler object in order to fill empty space that is not filled up by events.
func filler(events map[string]*Event) (f *Filler) {
	f = new(Filler)
	var bigCount, smallCount int
	for _, event := range events {
		switch event.Importance {
		case 1:
			bigCount += 1
		case 2:
			smallCount += 1
		}
	}
	fmt.Println(bigCount, smallCount)
	switch {
	case smallCount%4 != 0:
		f.Importance = 2
		return
	case bigCount%2 != 0:
		f.Importance = 1
		return
	}
	f = nil
	return
}

func MainHandler(c http.ResponseWriter, r *http.Request) {
	var obj EventsObj
	events := getEvents()
	// TODO: This is a bit slow, i imagine.
	eventarray := make(EventArray, len(events))
	i := 0
	for _, event := range events {
		eventarray[i] = event
		i++
	}
	sort.Sort(eventarray)
	obj.Events = eventarray
	obj.Filler = filler(events)
	// 	var p tmpl.PageInfo
	// 	p.Name = "events/main"
	// 	p.Title = "Events"
	// 	p.Object = obj
	// 	p.Request = r
	// 	p.Perms = auth.GetPerms(r)
	// 	tmpl.Execute(c, &p)
	tmpl.Render(c, r, "Events", "main", obj)
}

func EventHandler(c http.ResponseWriter, r *http.Request) {
	//var p tmpl.PageInfo
	//p.Name = "events/event"
	//p.Perms = auth.GetPerms(r)
	//p.Request = r

	id := r.URL.Path[len("/events/"):]
	fmt.Println(id)
	if len(id) == 0 {
		MainHandler(c, r)
		return
	}

	var event Event
	event.ID = id
	e := event.Load()
	if e != nil {
		//p.Name = "errors/500"
		fmt.Println(e)
	}

	e = event.LoadPage(r)
	if e != nil {
		//p.Name = "errors/500"
		fmt.Println(e)
	}

	//p.Title = event.Title
	//p.Object = event
	tmpl.Render(c, r, "Events - "+event.Title, "event", event)
}

// func ImgHandler(c http.ResponseWriter, r *http.Request) {
// 	path := "data/events/uploads/"+r.URL.Path[len("/events/img/"):]
// 	fmt.Println(path)
// 	http.ServeFile(c, r, path)
// }

func main() {
	fmt.Printf("Loading events server...\n")
	log.Println("testing!!!!")
	tmpl.SetModuleName("events")
	// Make required data directories
	os.MkdirAll("data/events/event", 0755)
	os.Mkdir("data/events/uploads", 0755)
	os.MkdirAll("data/pages/page/events", 0775)
	os.MkdirAll("data/pages/title/events", 0775)
	// Register handlers.
	http.HandleFunc("/events", MainHandler)
	http.HandleFunc("/events/", EventHandler)
	http.HandleFunc("/events/add", EditHandler)
	http.HandleFunc("/events/edit/", EditHandler)
	// "events/edit" POST goes to SaveHandler
	// 	http.HandleFunc("/events/img/", ImgHandler)

	http.ListenAndServe(":8130", nil)
}
