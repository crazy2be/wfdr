package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"
	// Local imports
	//"util/auth"
	"util/pages"
	tmpl "util/template"
)

const timeFormat = "2006-01-02 15:04 (Monday)"

func EditHandler(c http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		SaveHandler(c, r)
		return
	}
	//p := auth.GetPerms(r)
	//var i tmpl.PageInfo
	//i.Name = "events/edit"
	//i.Title = "Add Event"
	//i.Perms = p
	//i.Request = r

	// Makes the current data available to the template if available
	evEdit := len("/events/edit/")
	title := "Add Event"
	event := new(Event)
	if len(r.URL.Path) > evEdit {
		id := r.URL.Path[evEdit:]
		event.ID = id
		e := event.Load()
		if e != nil {
			fmt.Println(e)
		}
		e = event.LoadPage(r)
		if e != nil {
			fmt.Println(e)
		}
		title = "Edit Event - " + event.Title
	}

	tmpl.Render(c, r, title, "edit", event)
}

func getExt(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	default:
		return ""
	}
	panic("Not Reached!")
}

func saveFile(mpart *multipart.Part, id string) string {
	contentType := mpart.Header.Get("Content-Type")
	ext := getExt(contentType)

	if id == "" {
		// WARNING: Chance of a name collision, althrough extremely unlikely.
		id = strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	fileName := id + ext

	// WARNING: Should check file SIZE here...
	file, e := os.Create("data/img/base/" + fileName)
	if e != nil {
		fmt.Println(e)
		return ""
	}

	io.Copy(file, mpart)

	file.Close()
	return fileName
}

func HandleNess(c http.ResponseWriter, r *http.Request) (title, id, origID, date, desc, body []byte, fileName string, importance int, delete bool) {
	mreader, e := r.MultipartReader()
	if e != nil {
		return
	}
	for {
		mpart, e := mreader.NextPart()
		if mpart == nil {
			return
		}
		switch mpart.FormName() {
		case "title":
			title, e = ioutil.ReadAll(mpart)
		case "id":
			id, e = ioutil.ReadAll(mpart)
		case "orig-id":
			origID, e = ioutil.ReadAll(mpart)
		case "date":
			date, e = ioutil.ReadAll(mpart)
		case "desc":
			desc, e = ioutil.ReadAll(mpart)
		case "body":
			body, e = ioutil.ReadAll(mpart)
		case "delete-event":
			delete = true
			fmt.Println("!!!!Deleting event!!!!")
		case "importance":
			var imp []byte
			imp, e = ioutil.ReadAll(mpart)
			importance, e = strconv.Atoi(string(imp))
			fmt.Println("Importance:", string(imp), imp, importance)
		case "img":
			fileName = saveFile(mpart, string(id))
		}
		if e != nil {
			fmt.Println("Error:", e)
		}
	}
	return
}

func SaveHandler(c http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
	title, id, origID, dateString, desc, body, fileName, importance, delete := HandleNess(c, r)

	date, e := time.Parse(timeFormat, string(dateString))

	event := new(Event)

	// The ID changed, remove the old one.
	if !bytes.Equal(id, origID) && len(origID) > 0 {
		fmt.Println("ID Changed!")
		origEvent := new(Event)
		origEvent.ID = string(origID)
		origEvent.Load()

		// Copy the info to the new event
		*event = *origEvent

		origEvent.Delete()
	}

	event.ID = string(id)
	event.Load()

	// TODO: Delete the uploaded image as well!
	if delete {
		event.Delete()
		http.Redirect(c, r, "/events", 301)
		return
	}

	event.Importance = importance
	event.Title = string(title)
	event.Desc = string(desc)
	event.Link = "/events/" + string(id)
	event.Time = date
	if len(fileName) > 0 {
		event.Img = "/img/events/data/" + fileName
	}

	e = pages.SavePage("events/"+string(id), body, title)

	event.Save()

	fmt.Println(title, id, dateString, date, e, desc, body)
	fmt.Println(r.FormValue("img"))

	http.Redirect(c, r, "/events/"+event.ID, 301)
}
