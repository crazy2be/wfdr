// IN PROGRESS, INCOMPLETE!
// A simple authentication module that supports user and group based authentication methods, with users authenticated via e-mail addresses.
// In order to use it, you pass a http.Request to the Get() function, which returns the permissions of the current user based on the most permissive interpretation of their permissions and the permissions of their group.
// Uses the session library key "openid-email" to get the name of the user we are currently serving.
package perms

import (
	"http"
	"io/ioutil"
	"bytes"
	"fmt"
	"sort"
	// Local imports
	"github.com/crazy2be/session"
)

// TODO: Do we want more permissions here?
type Permissions struct {
	Read bool
	Write bool
	// Is the user authenticated at all? Aka 1) do they have a session cookie, and 2) have they logged in?
	Authenticated bool
	// Is the user recognized by the system? Do they have an account?
	Recognized bool
	// What path are these permissions for?
	Path string
}

type User struct {
	Email string
	Groups []string
	Perms []string
}

// Allows sorting the list of paths for matching.
type PermissionsList []Permissions

func (this PermissionsList) Len() int { 
	return len(this)
}
func (this PermissionsList) Less(i, j int) bool {
	return len(this[i].Path) > len(this[j].Path)
}
func (this PermissionsList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// DEPRECATED! Exists only for compatibility reasons. Can be removed when modules are changed to use the new function.
func GetPerms(r *http.Request) (p *Permissions) {
	return Get(r)
}

// Basic function that retrieves the permissions a user has based on the contents of their request, including cookies and request path. Designed to be a simple function for most uses. If you want more control, you can use the GetGroupPerms and GetUserPerms functions.
func Get(r *http.Request)  (p *Permissions) {
	p = new(Permissions)
	s, e := session.GetExisting(r)
	if e != nil {
		p.Authenticated = false
		return
	}
	p.Authenticated = true
	// Current authentication is based on e-mail. Might change this?
	uname := s.Get("openid-email")
	fmt.Println("Getting permissions for", uname)
	uperms := GetUserPerms(uname, r.URL.Path)
	if uperms == nil {
		p.Recognized = false
		return
	}
	p.Write = uperms.Write
	p.Read = uperms.Read
	fmt.Println("Grabbed permissions for user")
	groups := loadGroups(uname)
	for _, group := range(groups) {
		gperms := GetGroupPerms(group, r.URL.Path)
		if gperms == nil {
			continue
		}
		// Use the most permissive interpretation of the permissions. If a group is allowed to access something, so should all the users in the group.
		if !uperms.Read {
			if gperms.Read {
				p.Read = true
			}
		}
		if !uperms.Write {
			if gperms.Write {
				p.Write = true
			}
		}
	}
	return
}

// Retrieves the permissions for all members of a group with the given name.
func GetGroupPerms(name, path string) (p *Permissions) {
	mperms := loadPerms("data/shared/groups/"+name+"/perms")
	fmt.Println("Permissions for group", name, ":", mperms)
	p = matchPerms(mperms, path)
	fmt.Println("Matched permissions for group", name, ":", p)
	return
}

// Retrieves the user permissions, and the user permissions ONLY, for a user with a given name. Does not take group membership into account, and is likely not that useful for this reason.
func GetUserPerms(name, path string) (p *Permissions) {
	mperms := loadPerms("data/shared/users/"+name+"/perms")
	fmt.Println("Permissions for user", name, ":", mperms)
	p = matchPerms(mperms, path)
	fmt.Println("Matched permissions for user", name, ":", p)
	return
}

// Takes a list of permissions, likely garnered from a file, and returns the first found match. You should use sort.Sort to sort the list longest to shortest first, as this would allow the patterns to work as one would expect (more specific patterns override less specific ones, even if they have less permissions than than the more general ones).
func matchPerms(mperms PermissionsList, path string) (p *Permissions) {
	for _, mperm := range(mperms) {
		if pathMatch(mperm.Path, path) {
			// TODO: Caching
			fmt.Println("Found match for path", path, ":", mperm.Path)
			p = &mperm
			return
		}
	}
	return
}

func loadGroups(name string) (gr []string) {
	path := "data/shared/users/"+name+"/groups"
	file, e := ioutil.ReadFile(path)
	if e != nil {
		fmt.Println("Could not get group list for", name, ":", e)
	}
	lines := bytes.Split(file, []byte{'\n'}, -1)
	gr = make([]string, len(lines))
	// Is this necessary? It's somewhat inefficient...
	for i, line := range(lines) {
		gr[i] = string(line)
	}
	return
}

func loadPerms(path string) (mperms PermissionsList) {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		fmt.Println("Could not get group permissions for", path, ":", e)
		return
	}
	lines := bytes.Split(file, []byte("\n"), -1)
	mperms = make(PermissionsList, len(lines))
	for i, line := range(lines) {
		parts := bytes.Split(line, []byte(" "), 2)
		perms := mperms[i]
		for _, perm := range(parts[0]) {
			switch (perm) {
				case 'r':
					perms.Read = true
				case 'w':
					perms.Write = true
				default:
					fmt.Println("WARNING: Unrecognized permission", perm)
			}
			perms.Path = string(parts[1])
			mperms[i] = perms
		}
	}
	sort.Sort(mperms)
	if !sort.IsSorted(mperms) {
		fmt.Println("Failed to sort!")
	}
	return
}

// Does path match pattern?
// Stolen from HTTP library
func pathMatch(pattern, path string) bool {
	if len(pattern) == 0 {
		// should not happen
		return false
	}
	n := len(pattern)
	if pattern[n-1] != '/' {
		return pattern == path
	}
	return len(path) >= n && path[0:n] == pattern
}
