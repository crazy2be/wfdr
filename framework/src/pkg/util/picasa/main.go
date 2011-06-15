package picasa

import (
	"bufio"
	"http"
	"time"
	"net"
	"fmt"
	"os"
)

// The data structure used externally, including by templates. Should be succinct,
// only having fields that are likely needed externally. More can be added as needed.
type PhotoGroup struct {
	Date   time.Time
	Photos []Photo
}
type Photo struct {
	Title      string
	Summary    string
	Url        string // Image URL
	Width      int
	Height     int
	Modified   time.Time
	Published  *time.Time
	Thumbnails struct {
		Small, Medium, Large struct {
			Url    string
			Height int
			Width  int
		}
	}
}
type Album struct {
	Title     string
	Summary   string
	AlbumId   string
	Link      string // Album URL
	Url       string // Cover Image URL
	Modified  time.Time
	Published time.Time
}

// HTTP utility functions that really don't belong in this module, but are here until some other module needs them or someone moves them..

func NewRequest(url, authToken, method string) (req *http.Request) {
	req = new(http.Request)
	req.RawURL = url
	req.URL, _ = http.ParseURL(req.RawURL)
	req.Method = method
	req.Header = make(map[string][]string)
	// Set the authorization header with the token (required for picasa authentication).
	req.Header.Add("Authorization", "AuthSub  token="+authToken+"")
	return
}

func Send(req *http.Request) (resp *http.Response, e os.Error) {
	conn, e := Dial(req.URL.Host)
	e = SendTo(req, conn)
	resp, e = ReadFrom(conn, req.Method)
	conn.Close()
	return
}

func Dial(host string) (conn net.Conn, e os.Error) {
	// Open a connection
	conn, e = net.Dial("tcp", host+":80")
	if e != nil {
		fmt.Println("Error dialing host:", e)
	}
	return
}

func SendTo(req *http.Request, conn net.Conn) (e os.Error) {
	// Write our request struct to the connection in http wire format.
	e = req.Write(conn)
	if e != nil {
		fmt.Println("Error writing request:", e)
	}
	fmt.Printf("Wrote request\n")
	return
}

func ReadFrom(conn net.Conn, method string) (resp *http.Response, e os.Error) {
	// Read from and proccess the connection
	req := new(http.Request)
	req.Method = method
	
	reader := bufio.NewReader(conn)
	resp, e = http.ReadResponse(reader, req)
	if e != nil {
		fmt.Println("Error reading response:", e)
	}
	return
}
