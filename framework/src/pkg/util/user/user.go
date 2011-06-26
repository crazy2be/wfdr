// A simple package that allows persistent server-side storage of user settings and data.
package user

import (
	"http"
	"os"
	// Local imports
	"github.com/crazy2be/ini"
	"github.com/crazy2be/session"
)

type User struct {
	ID string // The user's ID, represented as an email in the current system.
	settings map[string]string // User settings, represented as a map[string]string right now. Might be represented as an interface{} later.
}

func Get(c http.ResponseWriter, r *http.Request) (u *User, err os.Error) {
	s := session.Get(c, r)
	if err != nil {
		return nil, err
	}
	u = new(User)
	u.ID = s.Get("openid-email")
	err = u.Load()
	return
}

func GetExisting(r *http.Request) (u *User, err os.Error) {
	s, err := session.GetExisting(r)
	if err != nil {
		return nil, err
	}
	u = new(User)
	u.ID = s.Get("openid-email")
	err = u.Load()
	return
}

func (u *User) Get(key string) string {
	return u.settings[key]
}

func (u *User) Set(key string, val string) {
	u.settings[key] = val
	u.Save()
}

// Where are the user's settings stored?
func (u *User) fileName() string {
	return "data/shared/users/"+u.ID+"/settings"
}

func (u *User) Load() (err os.Error) {
	filename := u.fileName()
	u.settings, err = ini.Load(filename)
	return err
}

func (u *User) Save() os.Error {
	filename := u.fileName()
	err := ini.Save(filename, u.settings)
	return err
}