package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	// Local imports
	"github.com/crazy2be/browser"
	"github.com/crazy2be/httputil"
	tmpl "util/template"
)

var Servers = make(map[string]string, 10)

// Map operations are not atomic in go.
var ServerLock sync.Mutex
var FileDirs = []string{"js", "css", "img"}

// URLs to allow through the firewall without authentication. Must match *exactly*!
var WhiteList = []string{"/pages/tos"}

func main() {
	fmt.Printf("Loading proxy server...\n")
	go runLogger()
	// TODO: Move this to external conf file.
	Servers = map[string]string{
		"photos": ":8090",
		"picasa": ":8090",
		"events": ":8130",
		// Main module (front page)
		"":       ":8140",
		"pages":  ":8150",
		"openid": ":8160",
		"news":   ":8170",
	}
	listener, e := net.Listen("tcp", ":8080")
	if e != nil {
		fmt.Println("Error opening listening port:", e)
		os.Exit(1)
	}
	// The listen-and-serve loop.
	for {
		proxyconn, e := listener.Accept()
		if e != nil {
			fmt.Println("Error accepting connection:", e)
			continue
		}
		//fmt.Println("[Accepted connection]")
		go func(conn net.Conn) {
			//defer fmt.Println("[Connection closed]")
			defer conn.Close()
			bufreader := bufio.NewReader(conn)
			for {
				startTime := time.Now().UnixNano()
				info, e := handleRequest(conn, bufreader)
				if e != nil {
					return
				}
				endTime := time.Now().UnixNano()
				deltaTime := float64(endTime-startTime) / (1000 * 1000 * 1000)

				info.TotalTime = deltaTime
				//fmt.Println(info)
				logChan <- info
				//fmt.Println("[Onto next http request...]")
			}
		}(proxyconn)
	}
}

func getFilePath(r *http.Request) string {
	path := r.URL.Path[1:]
	if browser.IsMobile(r) {
		pathsec := strings.SplitN(path, "/", 2)

		if len(pathsec) < 2 {
			return path
		}

		mobpath := "cache/" + pathsec[0] + "/mobile/" + pathsec[1]
		f, err := os.Open(mobpath)
		if err == nil {
			f.Close()
			return mobpath
		}
	} else {
		pathsec := strings.SplitN(path, "/", 2)

		if len(pathsec) < 2 {
			return path
		}

		deskpath := "cache/" + pathsec[0] + "/desktop/" + pathsec[1]
		f, e := os.Open(deskpath)

		if e == nil {
			f.Close()
			return deskpath
		}

	}
	return path
}

func handleRequest(c net.Conn, bufreader *bufio.Reader) (info *logInfo, e error) {
	req, e := http.ReadRequest(bufreader)
	if e != nil {
		fmt.Println("Error reading HTTP request,", e)
		return
	}
	// Logging stuff
	info = new(logInfo)
	info.Path = req.URL.Path
	info.Source = c.RemoteAddr().String()
	info.Referrer = req.Referer()
	info.ReqTime = time.Now().UnixNano()
	// The first directory, Split splits it a request for /css/main.css into "", "css", and "main.css".
	SubURL := strings.Split(req.URL.Path, "/")[1]
	// Check if the request URL matches one of the defined "file" URLs on the system. If it does, serve right away without any additional checks (for speed).
	if SubURL == "favicon.ico" {
		respwr := httputil.NewHttpResponseWriter(c)
		httputil.ServeFileOnly(respwr, req, "img/favicon.ico")
	}
	for _, path := range FileDirs {
		if path == SubURL {
			path := getFilePath(req)
			respwr := httputil.NewHttpResponseWriter(c)
			httputil.ServeFileOnly(respwr, req, path)
			return
		}
	}
	// Check if the user has the necessary authentication.
	if checkAuth(req) != true {
		if !isWhiteListed(req) {
			SubURL = "auth"
		}
	}
	info.Extra = fmt.Sprintln("Deffering to module:", SubURL)
	// 	fmt.Println("Time:", time.Seconds())
	// 	//fmt.Println("[Accepted HTTP request "+req.URL.Path+"]")
	// 	fmt.Println("Referer:", req.Referer)
	// 	fmt.Println("Source:", c.RemoteAddr())
	// Forward the request to another server that will actually handle it.
	ServerLock.Lock()
	addr, ok := Servers[SubURL]
	ServerLock.Unlock()
	if ok {
		fmt.Println("About to dial server")
		dialServer(c, req, addr)
	} else {
		error404(c, req)
	}
	c.Close()
	return
}

// Could be customized in order to require a password in a cookie or some other form of authentication, if desired. Unauthorized users will end up at the auth module, but still see the same URL on the front-end as they would normally. Thus, when they actually log in, they should see the proper page that they requested.
func checkAuth(req *http.Request) (Authenticated bool) {
	return true
}

func isWhiteListed(req *http.Request) bool {
	for _, url := range WhiteList {
		if req.URL.Path == url {
			return true
		}
	}
	return false
}

func dialServer(c net.Conn, r *http.Request, addr string) {
	// "localhost" can resolve to both the IPv4 and IPv6 addresses. Since we still use IPv4, we specify the address manually.
	addr = "127.0.0.1" + addr
	fmt.Println("Dialing remote server", addr, "in response to request for", r.URL.Path)
	serverconn, e := net.Dial("tcp", addr)
	if serverconn != nil {
		defer serverconn.Close()
	}
	if e != nil {
		error503(c, r, e)
		return
	}
	// Change the http version of the request so that the server doesn't think we support keepalive.
	// It seems to be the only way i can get it to work.
	r.Proto = "HTTP/1.0"
	r.ProtoMinor = 0
	r.Header.Del("Connection")

	fmt.Println("Writing data to server")
	r.Write(serverconn)
	serverconnbuf := bufio.NewReader(serverconn)
	fmt.Println("Reading response")

	resp, e := http.ReadResponse(serverconnbuf, r)
	if e != nil {
		error503(c, r, e)
		return
	}
	ct := resp.Header.Get("Content-Type")
	fmt.Println("Content-Type of response:", ct)
	if strings.Contains(ct, "text/xml") {
		resp.Header.Set("Content-Type", "text/html; charset=utf-8")
	}
	//resp.AddHeader("Connection", "close")
	//resp.ProtoMinor = 0
	fmt.Println("Writing response")
	resp.Write(c)
	resp.Body.Close()
	fmt.Println("Done!")
	//c.Close()
	return
	//conn.Close()
	//io.Copy(conn, serverconn)
	//serverconn.Close()

}

// TODO: These should really be moved to external files, perhaps mustache templates.
func error503(c net.Conn, r *http.Request, e error) {
	fmt.Println("Error dialing server:", e)
	httpError(c, r, "503", "Service Unavailable")
}

func error404(c net.Conn, r *http.Request) {
	fmt.Println("404 Error")
	httpError(c, r, "404", "Not Found")
}

func httpError(c net.Conn, r *http.Request, code string, desc string) {
	var i tmpl.PageInfo
	i.Name = "shared/errors/" + code
	i.Title = "Error " + code + ": " + desc
	i.Request = r
	buf := bytes.NewBuffer([]byte(""))
	tmpl.Execute(buf, &i)
	fmt.Fprintln(c, "HTTP/1.1", code, desc)
	fmt.Fprintln(c, "Content-Type:", "text/html")
	fmt.Fprintln(c, "Content-Length:", buf.Len())
	fmt.Fprintln(c, "")
	io.Copy(c, buf)
}

// Infomation that should be logged for every HTTP request that passes through this proxy server.
type logInfo struct {
	Path     string
	Source   string
	Referrer string
	// Time that the request was sent, in nanoseconds.
	ReqTime int64
	// Extra info that only applies to specific requests.
	Extra string
	// Render/send time
	TotalTime float64
}

// Is 100 a reasonable number?
var logChan = make(chan *logInfo, 100)

func runLogger() {
	// Make logging directory
	os.MkdirAll("log/server", 0755)

	localt := time.Now()
	lfname := localt.Format("2006-01-02 15:04:05")

	lfile, err := os.Create("log/server/" + lfname)
	if err != nil {
		log.Fatal("Failed to open log! Aborting because you should really have a log set up, and something is probably horribly broken...")
	}
	defer lfile.Close()

	lenc := json.NewEncoder(lfile)

	for {
		linfo := <-logChan
		err := lenc.Encode(linfo)
		if err != nil {
			fmt.Println("Warning: Failed to log request!", err)
		}
	}
}
