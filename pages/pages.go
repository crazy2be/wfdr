package pages

import (
	"io/ioutil"
	"net/http"
	"strings"
	"bytes"
	"path"
	"log"
	"os"

	tmpl "wfdr/template"
)

const (
	TITLE_DIRECTORY   = "data/pages/title/"
	CONTENT_DIRECTORY = "data/tmpl/base/page/"
)

type PageContent []byte

func (pc PageContent) String() string {
	return string([]byte(pc))
}

type PageData struct {
	Content PageContent
	Title   string
	Name    string
}

func LoadPage(name string) (p *PageData, e error) {
	p = new(PageData)
	var page bytes.Buffer

	// Prevents users from sticking ../../ etc in URLs and writing to other files.
	// TODO: path.Clean()?
	if strings.Index(name, ".") != -1 {
		p.Content = PageContent([]byte("Nice try ;)."))
		p.Title = "HAX AVOIDED"
		return
	}

	titleb, e := ioutil.ReadFile(TITLE_DIRECTORY + name)
	title := string(titleb)
	if e != nil {
		log.Println("Oops, could not load title for page", name, ". Error:", e)
		title = "[No Title]"
	}

	inf := &tmpl.PageInfo{Title: title, Name: "data/page/"+name}
	tmpl.Execute(&page, inf)
	
	p.Content = page.Bytes()
	p.Title = title
	p.Name = name
	return
}

// Saves the page to a file in data/tmpl/base/page, with the title in data/pages/title.
func (p *PageData) Save() (e error) {
	log.Println(p.Content, p.Title, p.Name)

	nameDir, _ := path.Split(p.Name)

	os.MkdirAll(TITLE_DIRECTORY+nameDir, 0755)
	os.MkdirAll(CONTENT_DIRECTORY+nameDir, 0755)

	e = ioutil.WriteFile(CONTENT_DIRECTORY+p.Name, p.Content, 0666)
	e = ioutil.WriteFile(TITLE_DIRECTORY+p.Name, []byte(p.Title), 0666)
	return
}

// Deletes the page (from disk!).
func (p *PageData) Delete() (err error) {
	err = os.Remove(TITLE_DIRECTORY + p.Name)
	if err != nil {
		return
	}
	err = os.Remove(CONTENT_DIRECTORY + p.Name)
	if err != nil {
		return
	}
	return
}

// DEPRECATED! USE PageData.Save() instead!
func SavePage(pageName string, content, title []byte) (e error) {
	p := &PageData{content, string(title), pageName}
	e = p.Save()
	return
}

// DEPRECATED! Use LoadPage() instead!
func GetPageData(pageName string, r *http.Request) (p *PageData, e error) {
	p, e = LoadPage(pageName)
	return
}

type Page struct {
	Name  string
	Title string
}

func GetPageList() (pages []Page) {
	pagenames := getPageListInDirectory("")
	if pagenames == nil {
		return
	}

	for _, pagename := range pagenames {
		var page Page
		page.Name = pagename
		title, err := ioutil.ReadFile(TITLE_DIRECTORY + pagename)
		if err != nil {
			log.Println("Error getting page list:", err)
			continue
		}
		page.Title = string(title)
		pages = append(pages, page)
	}
	return
}

func getPageListInDirectory(dir string) (pages []string) {
	pagesFolder, err := os.Open(TITLE_DIRECTORY + dir)
	if err != nil {
		log.Println(err)
		return
	}

	pagesinfo, err := pagesFolder.Readdir(-1)
	if err != nil {
		log.Println(err)
		return
	}

	for _, pageinfo := range pagesinfo {
		if pageinfo.IsDir() {
			pages = append(pages, getPageListInDirectory(dir+"/"+pageinfo.Name())...)
		} else {
			pages = append(pages, path.Join(dir, pageinfo.Name()))
		}
	}
	return
}
