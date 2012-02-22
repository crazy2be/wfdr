package pages

import (
	"net/http"
	"path"
	
	"wfdr/perms"
	"wfdr/template"
)

type PageServer struct {
	Prefix    string
	PageAlias string // "Article" for news module, etc.
	Manager   Manager
}

func (ps *PageServer) Save(c http.ResponseWriter, r *http.Request) {
	oldname := r.FormValue("oldname")
	name := r.FormValue("name")
	
	if !perms.ToEditPage(r, path.Join(ps.Prefix, name)) {
		template.Error403(c, r, name)
		return
	}
	
	if !perms.ToEditPage(r, path.Join(ps.Prefix, oldname)) {
		template.Error403(c, r, oldname)
		return
	}
	
	content := r.FormValue("content")
	title := r.FormValue("title")
	
	err := ps.Manager.Save(name, title, []byte(content))
	if err != nil {
		template.Error500(c, r, err)
		return
	}
	
	if oldname != name && oldname != "" {
		err := ps.Manager.Delete(oldname)
		if err != nil {
			template.Error500(c, r, err)
			return
		}
	}
	
	http.Redirect(c, r, path.Join(ps.Prefix, name), 301)
}

func (ps *PageServer) Display(c http.ResponseWriter, r *http.Request) {
	page, err := ps.Manager.Load(r.URL.Path[len(ps.Prefix):])
	if err != nil {
		template.Error404(c, r, err)
		return
	}
	template.Render(c, r, page.Title, "page", page)
}

func (ps *PageServer) Edit(c http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) < len(ps.Prefix) + len("/edit") {
		template.Error404(c, r, nil)
		return
	}
	
	name := r.URL.Path[len(ps.Prefix) : len(r.URL.Path)-len("/edit")]
	
	page, err := ps.Manager.Load(name)
	if err != nil {
		template.Error500(c, r, err)
		return
	}
	
	template.Render(c, r, "Editing "+page.Title, "edit", page)
}

func (ps *PageServer) Add(c http.ResponseWriter, r *http.Request) {
	var page Page
	page.Title = "[New "+ps.PageAlias+"]"
	template.Render(c, r, "Add New "+ps.PageAlias, "edit", page)
}

func (ps *PageServer) List(c http.ResponseWriter, r *http.Request) {
	plist, err := ps.Manager.List()
	if err != nil {
		template.Error500(c, r, err)
		return
	}
	
	template.Render(c, r, ps.PageAlias, "main", plist)
	return
}