package browser
// A simple package to implement browser checking to allow mobile-specific, browser specific (Helllllo IE!), and platform specific logic.

import (
	"strings"
	"http"
	"fmt"
)

var mobileAgentSubstrings = [...]string{"iPod", "iPhone", "Mobile", "Phone", "Android"}

func IsMobile(r *http.Request) bool {
	agent := r.UserAgent
	
	for _, substr := range mobileAgentSubstrings {
		if strings.Contains(agent, substr) {
			fmt.Println("Mobile browser detected! User-Agent:", agent)
			return true
		}
	}
	return false
}
