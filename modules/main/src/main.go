package main

import (
	"net/http"
	// Local imports
	"util/template"
)

func Handler(c http.ResponseWriter, r *http.Request) {
	template.Render(c, r, "Main", "main", nil)
}

func main() {
	http.HandleFunc("/", Handler)
	http.ListenAndServe(":8140", nil)
}
