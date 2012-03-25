// A simple package that allows persistent server-side storage of user settings and data.
package user

import (
	"net/http"

	"wfdr/session"

	"github.com/crazy2be/ini"
)

type User struct {
	ID       string            // The user's ID, represented as an email in the current system.
	settings map[string]string // User settings, represented as a map[string]string right now. Might be represented as an interface{} later.
}

func Get(r *http.Request) (*User, error) {
	var err error
	u := new(User)
	u.ID, err = session.Get(r, "openid-email")
	if err != nil {
		return nil, err
	}
	err = u.Load()
	if err != nil {
		return nil, err
	}
	return u, nil
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
	return "data/shared/users/" + u.ID + "/settings"
}

func (u *User) Load() (err error) {
	filename := u.fileName()
	u.settings, err = ini.Load(filename)
	return err
}

func (u *User) Save() error {
	filename := u.fileName()
	err := ini.Save(filename, u.settings)
	return err
}
