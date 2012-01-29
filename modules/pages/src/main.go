package main

import (
	"fmt"
	"net/http"
	"strings"
	// Local imports
	"github.com/crazy2be/perms"
	"util/pages"
	tmpl "util/template"
)

func SaveHandler(c http.ResponseWriter, r *http.Request) {
	p := perms.GetPerms(r)
	if !p.Write {
		fmt.Fprintln(c, "Access Denied")
		return
	}

	pageName := r.URL.Path[len("/pages/"):]

	if strings.Index(pageName, ".") != -1 {
		return
	}

	content := []byte(r.FormValue("content"))
	title := []byte(r.FormValue("title"))

	e := pages.SavePage(pageName, content, title)
	if e != nil {
		fmt.Fprintln(c, e)
		return
	}
	http.Redirect(c, r, "/pages/"+pageName, 301)
}
func EditHandler(c http.ResponseWriter, r *http.Request) {
	//var p tmpl.PageInfo
	//p.Name = "pages/edit"
	//p.Request = r

	//perms := auth.GetPerms(r)
	//p.Perms = perms

	pagePath := r.URL.Path[len("/pages/") : len(r.URL.Path)-len("/edit")]
	fmt.Println(pagePath)
	pageData, _ := pages.GetPageData(pagePath, r)
	//if e != nil {
	//	p.Name = "errors/404"
	//}

	//p.Object = pageData

	fmt.Println("Request for pages server. Responding.")
	tmpl.Render(c, r, "Editinng "+pageData.Title, "edit", pageData)
}

func ListHandler(c http.ResponseWriter, r *http.Request) {
	plist := pages.GetPageList()
	tmpl.Render(c, r, "Pages list", "list", plist)
	return
}

func Handler(c http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		SaveHandler(c, r)
		return
	}
	rp := r.URL.Path
	fmt.Println(rp, rp[len(rp)-len("/edit"):])
	if rp[len(rp)-len("/edit"):] == "/edit" {
		EditHandler(c, r)
		return
	}
	if len(rp) == len("/pages/") {
		ListHandler(c, r)
		return
	}

	pageData, e := pages.GetPageData(r.URL.Path[len("/pages/"):], r)
	if e != nil {
		// TODO: 404 Error
		fmt.Println("Info: Page not found")
	}
	fmt.Println("Request for pages server. Responding.")
	//tmpl.Execute(c, &p)
	tmpl.Render(c, r, pageData.Title, "main", pageData)
}

func main() {
	fmt.Printf("Loading pages server...\n")
	tmpl.SetModuleName("pages")
	http.HandleFunc("/pages/", Handler)
	//http.HandleFunc("/pages/edit/", EditHandler)
	//http.HandleFunc("/pages/save/", SaveHandler)
	http.ListenAndServe(":8150", nil)
}
