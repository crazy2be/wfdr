package pages

type Info struct {
	Title string
	Name  string
}

type Content []byte

func (pc Content) String() string {
	return string([]byte(pc))
}

type Page struct {
	Info
	Content
}

var defaultManager *FSManager = &FSManager{"data/pages/title/", "data/pages/content/"}

func Load(name string) (p *Page, err error) {
	return defaultManager.Load(name)
}

func Save(name string, title string, content []byte) error {
	return defaultManager.Save(name, title, content)
}

func Move(oldname string, newname string) error {
	return defaultManager.Move(oldname, newname)
}

func Delete(name string) error {
	return defaultManager.Delete(name)
}

func List() ([]Info, error) {
	return defaultManager.List()
}