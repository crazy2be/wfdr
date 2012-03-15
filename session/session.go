// A simple package that allows persistent server-side storage of session settings through a variety of configurable backends. Typical usage is just 
//	err := session.Set(c, r, "foo", "bar")
//	val, err := session.Get(r, "foo")
package session

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/crazy2be/ini"
)

var defaultManager Manager

// GenID returns a new randomly generated session identifier. The optional argument allows the caller to inject additional entropy into the session generation process.
func GenID(entropy []byte) string {
	seed := time.Now().UnixNano()
	hash := sha256.New()

	fmt.Fprintf(hash, "%d", seed)
	if entropy != nil {
		fmt.Fprintf(hash, "%b", entropy)
	}

	id := make([]byte, 0)
	id = hash.Sum(id)
	return base64.URLEncoding.EncodeToString(id)
}

// Get attempts to find the sessionid cookie in r, then delegates to the defaultManager, which tries to find the value of the element corresponding to key. Returns an empty string and an error if the session or key does not exist.
func Get(r *http.Request, key string) (string, error) {
	return GetWith(defaultManager, r, key)
}

// GetWith is the same as Get(), but allows you to specify a manager different from the defaultManager.
func GetWith(m Manager, r *http.Request, key string) (string, error) {
	cookie, err := r.Cookie("sessionid")
	if err != nil {
		return "", errors.New("No sessionid found!" + err.Error())
	}

	return m.Get(cookie.Value, key)
}

// Set first attempts to find the current session associated with r, if any. If that fails, it will create a new session and associate it with c by setting the sessionid cookie. Finally, it will set key to val using the defaultManager. Note that this must be called before any data is sent on c in order to ensure that it has an effect if a session does not already exist.
func Set(c http.ResponseWriter, r *http.Request, key, val string) error {
	return SetWith(defaultManager, c, r, key, val)
}

// SetWith is the same as Set(), but allows you to specify a manager different from the defaultManager.
func SetWith(m Manager, c http.ResponseWriter, r *http.Request, key, val string) error {
	id := ""
	cookie, err := r.Cookie("sessionid")
	if err != nil {
		id = GenID([]byte(r.RemoteAddr))
		c.Header().Add("Cookie", "sessionid="+id+"; path=/")
	} else {
		id = cookie.Value
	}
	return m.Set(id, key, val)
}

// FSManager is the current default storage mechanism, relying on the local filesystem to store session information in ini files.
type FSManager struct {
	path string
}

// NewFSManager creates a new FSManager with storage in the location pointed at by dir, creating the directory if necessary. Returns an error if the directory pointed at by dir could not be created.
func NewFSManager(dir string) (*FSManager, error) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
	fsm := &FSManager{dir}
	return fsm, nil
}

// Get attempts to find the file corresponding to id, searches through it for key, and returns the corresponding value if it exists. Returns an empty string and an error if the specified id or key does not exist.
func (fsm *FSManager) Get(id, key string) (string, error) {
	filename := path.Join(fsm.path, id)
	vals, err := ini.Load(filename)
	if err != nil {
		return "", err
	}
	val, ok := vals[key]
	if !ok {
		return "", fmt.Errorf("Specified key '%s' does not exist in file '%s' (for session '%s').", key, filename, id)
	}
	return val, nil
}

// Set finds the file corresponding to id, loads it into memory, sets the corresponding key to the specified val, and then writes it back to disk. Returns an error, if any.
func (fsm *FSManager) Set(id, key, val string) error {
	filename := path.Join(fsm.path, id)
	vals, err := ini.Load(filename)
	if err != nil {
		// Assume no session file exists yet
		vals = make(map[string]string, 1)
	}
	vals[key] = val
	err = ini.Save(filename, vals)
	if err != nil {
		return err
	}
	return nil
}

// Manager describes a type which knows how to store and handle session information. It could be backed by a filesystem, database, or even just an in-memory cache depending on what makes the most sense for the application at hand.
type Manager interface {
	// Get retreives the value corresponding to the given id and key from whichever backend it corresponds to. Returns an error if either id or key could not be found in the backing datastore.
	Get(id, key string) (string, error)
	// Set sets the value corresponding to the given id and key, returning an error if the process was unsucessfull.
	Set(id, key, val string) error
}

func init() {
	var err error
	defaultManager, err = NewFSManager("data/shared/sessions")
	if err != nil {
		panic(err)
	}
}
