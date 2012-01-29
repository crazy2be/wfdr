package main

import (
	"fmt"
	"net/http"
	// Local imports
	"github.com/crazy2be/perms"
	"util/pages"
	"util/template"
	//"util/dlog"
)

func checkPerms(c http.ResponseWriter, r *http.Request, name string) bool {
	r.URL.Path = "/news/edit/" + name

	p := perms.GetPerms(r)
	if !p.Write {
		fmt.Fprintln(c, "Access Denied")
		return false
	}
	return true
}

func SaveHandler(c http.ResponseWriter, r *http.Request) {
	oldname := r.FormValue("oldname")
	name := r.FormValue("name")

	if !checkPerms(c, r, name) {
		return
	}

	// Make sure the user has permission to delete the old event as well.
	if oldname != name && len(oldname) > 0 {
		if !checkPerms(c, r, name) {
			return
		}
		oldpage, err := pages.LoadPage(oldname)
		if err == nil {
			oldpage.Delete()
		}
	}

	content := []byte(r.FormValue("content"))
	title := []byte(r.FormValue("title"))

	e := pages.SavePage(name, content, title)
	if e != nil {
		fmt.Fprintln(c, e)
		return
	}
	http.Redirect(c, r, "/news/"+name, 301)
}

func EditHandler(c http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) < len("/news//edit") {
		template.Error404(c, r, nil)
		return
	}
	pagePath := r.URL.Path[len("/news/") : len(r.URL.Path)-len("/edit")]
	fmt.Println(pagePath)
	pageData, _ := pages.GetPageData(pagePath, r)

	template.Render(c, r, "Editing "+pageData.Title, "edit", pageData)
}

func AddHandler(c http.ResponseWriter, r *http.Request) {
	var pageData pages.PageData
	pageData.Title = "[New Article]"
	template.Render(c, r, "Add News Article", "edit", pageData)
}

func ArticleHandler(c http.ResponseWriter, r *http.Request) {
	pageData, e := pages.GetPageData(r.URL.Path[len("/news/"):], r)
	if e != nil {
		fmt.Println("Info: Page not found")
		template.Error404(c, r, nil)
		return
	}
	template.Render(c, r, pageData.Title, "article", pageData)
}

func Handler(c http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		SaveHandler(c, r)
		return
	}
	rp := r.URL.Path
	fmt.Println(rp, rp[len(rp)-len("/edit"):])
	// /news/<article>/edit triggers this
	if rp[len(rp)-len("/edit"):] == "/edit" {
		EditHandler(c, r)
		return
	}
	if len(rp) != len("/news/") {
		ArticleHandler(c, r)
		return
	}

	plist := pages.GetPageList()
	template.Render(c, r, "News", "main", plist)
	return
}

func main() {
	fmt.Printf("Loading news server...\n")
	template.SetModuleName("news")
	http.HandleFunc("/news/add", AddHandler)
	http.HandleFunc("/news/", Handler)
	http.ListenAndServe(":8170", nil)
}
