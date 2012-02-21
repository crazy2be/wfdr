package pages

import (
	"io/ioutil"
	"net/http"
	"path"
	"log"
	"os"
)

const (
	TITLE_DIRECTORY   = "data/pages/title/"
	CONTENT_DIRECTORY = "data/pages/content/"
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

func LoadPage(name string) (p *PageData, err error) {
	p = new(PageData)
	p.Name = path.Clean(name)

	titleb, err := ioutil.ReadFile(path.Join(TITLE_DIRECTORY, p.Name))
	p.Title = string(titleb)
	if err != nil {
		p.Title = "Error reading page title: " + err.Error()
		return p, err
	}
	
	p.Content, err = ioutil.ReadFile(path.Join(CONTENT_DIRECTORY, p.Name))
	if err != nil {
		p.Content = []byte("Error reading page content: " + err.Error())
		return p, err
	}
	return p, nil
}

// Saves the page to a file in data/tmpl/base/page, with the title in data/pages/title.
func (p *PageData) Save() error {
	nameDir, _ := path.Split(p.Name)

	os.MkdirAll(path.Join(CONTENT_DIRECTORY, nameDir), 0755)
	os.MkdirAll(path.Join(TITLE_DIRECTORY, nameDir), 0755)

	err1 := ioutil.WriteFile(path.Join(CONTENT_DIRECTORY, p.Name), p.Content, 0666)
	err2 := ioutil.WriteFile(path.Join(TITLE_DIRECTORY, p.Name), []byte(p.Title), 0666)
	
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

// Deletes the page (from disk!).
func (p *PageData) Delete() error {
	err1 := os.Remove(path.Join(TITLE_DIRECTORY, p.Name))
	err2 := os.Remove(path.Join(CONTENT_DIRECTORY, p.Name))
	
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
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
