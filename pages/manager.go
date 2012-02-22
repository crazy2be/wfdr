package pages

import (
	"io/ioutil"
	"path"
	"os"
)

type Manager interface {
	Load(name string) (*Page, error)
	Save(name string, title string, content []byte) error
	Move(oldname string, newname string) error
	Delete(name string) error
	List() ([]Info, error)
}

type FSManager struct {
	TitleDir   string
	ContentDir string
}

func (fsm *FSManager) Load(name string) (*Page, error) {
	p := new(Page)
	p.Name = path.Clean(name)
	
	titleb, err1 := ioutil.ReadFile(path.Join(fsm.TitleDir, p.Name))
	p.Title = string(titleb)
	if err1 != nil {
		p.Title = "Error reading page title: " + err1.Error()
	}
	
	var err2 error
	p.Content, err2 = ioutil.ReadFile(path.Join(fsm.ContentDir, p.Name))
	if err2 != nil {
		p.Content = []byte("Error reading page content: " + err2.Error())
	}
	
	if err1 != nil {
		return p, err1
	} else if err2 != nil {
		return p, err2
	}
	return p, nil
}

func (fsm *FSManager) Save(name string, title string, content []byte) (err error) {
	dir, _ := path.Split(name)

	err = os.MkdirAll(path.Join(fsm.ContentDir, dir), 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(path.Join(fsm.TitleDir, dir), 0755)
	if err != nil {
		return err
	}

	err1 := ioutil.WriteFile(path.Join(fsm.ContentDir, name), content, 0666)
	err2 := ioutil.WriteFile(path.Join(fsm.TitleDir, name), []byte(title), 0666)
	
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

// Move moves the contents of the page at oldname to newname.
// The current implementation does this by loading the old page, saving it under the new name, and then deleting the old page.
// In the event of a failure, the worst case scenerio is that two copies of the page will exist.
func (fsm *FSManager) Move(oldname string, newname string) error {
	page, err := Load(oldname)
	if err != nil {
		return err
	}
	
	err = Save(newname, page.Title, page.Content)
	if err != nil {
		return err
	}
	
	err = Delete(oldname)
	if err != nil {
		return err
	}
	
	return nil
}

func (fsm *FSManager) Delete(name string) (error) {
	err1 := os.Remove(path.Join(fsm.TitleDir, name))
	err2 := os.Remove(path.Join(fsm.ContentDir, name))
	
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

func (fsm *FSManager) List() ([]Info, error) {
	pages := make([]Info, 0)
	names, err1 := fsm.listInDir("")
	if names == nil {
		return nil, err1
	}

	for _, name := range names {
		page := Info{Name: name}
		title, err := ioutil.ReadFile(path.Join(fsm.TitleDir, name))
		if err != nil {
			return pages, err1
		}
		page.Title = string(title)
		pages = append(pages, page)
	}
	return pages, nil
}

// listInDir recursively calls itself at each level of the filesystem tree, constructing a list of pages available in a []string. It returns the list of pages constructed thus far, and the error, if any.
func (fsm *FSManager) listInDir(dir string) ([]string, error) {
	pages := make([]string, 0)
	dirf, err := os.Open(path.Join(fsm.TitleDir, dir))
	if err != nil {
		return nil, err
	}
	
	pagesinfo, err := dirf.Readdir(-1)
	if err != nil {
		return nil, err
	}
	
	for _, fi := range pagesinfo {
		if fi.IsDir() {
			morepages, err := fsm.listInDir(path.Join(dir, fi.Name()))
			if err != nil {
				return pages, err
			}
			pages = append(pages, morepages...)
		} else {
			pages = append(pages, path.Join(dir, fi.Name()))
		}
	}
	return pages, nil
}