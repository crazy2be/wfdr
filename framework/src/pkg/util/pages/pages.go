package pages

import (
	"os"
	"fmt"
	"log"
	"http"
	"path"
	"bytes"
	"strings"
	"io/ioutil"
	// Local imports
	tmpl "util/template"
	"util/dlog"
)

const (
	TITLE_DIRECTORY = "data/pages/title/"
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

func LoadPage(name string) (p *PageData, e os.Error) {
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
		fmt.Println("Oops, could not load title for page", name, ". Error:", e)
		title = "[No Title]"
	}
	
	// TODO: 404 error if page is not loaded
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		log.Println(err)
	}
	tmpl.Render(&page, req, string(title), "data/page/"+name, nil)
	log.Println(&page)
	p.Content = page.Bytes()
	p.Title = title
	p.Name = name
	return
}

// Saves the page to a file in data/tmpl/base/page, with the title in data/pages/title.
func (p *PageData) Save() (e os.Error) {
	fmt.Println(p.Content, p.Title, p.Name)

	nameDir, _ := path.Split(p.Name)

	os.MkdirAll(TITLE_DIRECTORY + nameDir, 0755)
	os.MkdirAll(CONTENT_DIRECTORY + nameDir, 0755)

	e = ioutil.WriteFile(CONTENT_DIRECTORY + p.Name, p.Content, 0666)
	e = ioutil.WriteFile(TITLE_DIRECTORY + p.Name, []byte(p.Title), 0666)
	return
}

// Deletes the page (from disk!).
func (p *PageData) Delete() (err os.Error) {
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
func SavePage(pageName string, content, title []byte) (e os.Error) {
	p := &PageData{content, string(title), pageName}
	e = p.Save()
	return
}

// DEPRECATED! Use LoadPage() instead!
func GetPageData(pageName string, r *http.Request) (p *PageData, e os.Error) {
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
			dlog.Println(err)
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
		dlog.Println(err)
		return
	}

	pagesinfo, err := pagesFolder.Readdir(-1)
	if err != nil {
		dlog.Println(err)
		return
	}

	for _, pageinfo := range pagesinfo {
		if pageinfo.IsDirectory() {
			pages = append(pages, getPageListInDirectory(dir+"/"+pageinfo.Name)...)
		} else {
			pages = append(pages, path.Join(dir, pageinfo.Name))
		}
	}
	return
}
