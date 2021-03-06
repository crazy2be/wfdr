package picasa

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
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
	Published  time.Time
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

const TimeFormat = "2006-01-02T15:04:05"

// Parses a timestamp of the format 2011-05-25T12:50:49 (as returned by picasa) into a time.Time object.
func ParseTimestamp(timestamp string) (d time.Time, err error) {
	d, err = time.Parse(TimeFormat, timestamp)
	return
	// 	// Splits into time and date parts
	// 	var dateTime = strings.Split(timestamp, "T")
	// 	// Splits off the timezone
	// 	//var timeZone = strings.Split(dateTime[1], "Z", -1)
	// 	
	// 	// Parse the date
	// 	var date = strings.Split(dateTime[0], "-")
	// 	d.Year, _ = strconv.ParseInt(date[0], 10, 64)
	// 	d.Month, _ = strconv.Atoi(date[1])
	// 	d.Day, _ = strconv.Atoi(date[2])
	// 
	// 	// Parse the time
	// 	var time = strings.Split(dateTime[1], ":")
	// 	d.Hour, _ = strconv.Atoi(time[0])
	// 	d.Minute, _ = strconv.Atoi(time[1])
	// 	d.Second, _ = strconv.Atoi(time[2])
}

// HTTP utility functions that really don't belong in this module, but are here until some other module needs them or someone moves them..

func NewRequest(surl, authToken, method string) (req *http.Request) {
	req = new(http.Request)
	req.URL, _ = url.Parse(surl)
	req.Method = method
	req.Header = make(map[string][]string)
	// Set the authorization header with the token (required for picasa authentication).
	req.Header.Add("Authorization", "AuthSub  token="+authToken+"")
	return
}

func Send(req *http.Request) (resp *http.Response, e error) {
	conn, e := Dial(req.URL.Host)
	e = SendTo(req, conn)
	resp, e = ReadFrom(conn, req.Method)
	conn.Close()
	return
}

func Dial(host string) (conn net.Conn, e error) {
	// Open a connection
	conn, e = net.Dial("tcp", host+":80")
	if e != nil {
		fmt.Println("Error dialing host:", e)
	}
	return
}

func SendTo(req *http.Request, conn net.Conn) (e error) {
	// Write our request struct to the connection in http wire format.
	e = req.Write(conn)
	if e != nil {
		fmt.Println("Error writing request:", e)
	}
	fmt.Printf("Wrote request\n")
	return
}

func ReadFrom(conn net.Conn, method string) (resp *http.Response, e error) {
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
