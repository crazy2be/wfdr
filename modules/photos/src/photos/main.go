package main

import (
	"fmt"
	"net/http"
	// Local imports
	"util/template"
)

func main() {
	fmt.Printf("Loading photos server...\n")
	template.SetModuleName("photos")
	http.HandleFunc("/photos/", Handler)
	http.HandleFunc("/photos/upload", UploaderHandler)
	// TODO: Move these to their own module.
	// We have to get sessions to transfer between processes correctly before we can do that, however.
	http.HandleFunc("/picasa/auth", AuthHandler)
	http.HandleFunc("/picasa/upload", UploadHandler)
	http.ListenAndServe(":8090", nil)
}
