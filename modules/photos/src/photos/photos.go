package main

import (
	"fmt"
	"net/http"
	// Local imports
	"github.com/crazy2be/user"
	tmpl "util/template"
)

func Handler(c http.ResponseWriter, r *http.Request) {
	var photosLen = len("/photos/")
	// On the root
	if len(r.URL.Path) == photosLen {
		MainHandler(c, r)
		return
	} else {
		AlbumHandler(c, r)
		return
	}
}

func MainHandler(c http.ResponseWriter, r *http.Request) {
	var albums = GetAlbums()
	fmt.Printf("%d Albums\n", len(*albums))
	tmpl.Render(c, r, "Photos", "main", &albums)
}

func AlbumHandler(c http.ResponseWriter, r *http.Request) {
	var albumName = r.URL.Path[len("/photos/"):]
	fmt.Printf("Album ID requested: %s\n", albumName)
	var photos = GetPhotos(albumName)
	fmt.Printf("%d Photos\n", len(photos))
	tmpl.Render(c, r, "Photos - "+albumName, "album", &photos)
}

type UploaderData struct {
	Albums *[]Album
	// Is the user authenticated to use picasa? They have to give us access to their picasa account.
	PicasaAuthenticated bool
}

func UploaderHandler(c http.ResponseWriter, r *http.Request) {
	u, _ := user.Get(c, r)
	token := u.Get("picasa-authsub-token")
	fmt.Println("Host:", r.Host)
	albums := GetAlbums()

	data := new(UploaderData)
	data.Albums = albums
	data.PicasaAuthenticated = (token != "")
	// TODO: Make photos use the new login system.

	tmpl.Render(c, r, "Upload Photos", "uploader", data)
}
