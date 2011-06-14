package http

import (
	"http"
	"fmt"
	"io"
	"os"
	"time"
	"mime"
	pathm "path"
)

// Mostly depreciated, should not be needed for new code given that r.Cookie is now available. However, it's still useful to get a particular cookie until i write a function to do that...
func GetCookies(r *http.Request) map[string]string {
	//cookiestring := r.Header["Cookie"]
	//cookiearray := strings.Split(cookiestring, "; ", -1)
	cookiearray := r.Cookie
	fmt.Println(r.Cookie)
	//fmt.Println(r.Header["Cookie"])
	//cookiearray := r.Header["Cookie"]
	cookies := make(map[string]string, len(cookiearray))
	// Construct a map of cookie names and values.
	for _, cookie := range cookiearray {
		cookies[cookie.Name] = cookie.Value
		// 		cookienamevalue := strings.Split(string(cookiearray[cookie]), "=", 2)
		// 		if len(cookienamevalue) < 2 {
		// 			continue;
		// 		}
		// 		cookies[cookienamevalue[0]] = cookienamevalue[1]
	}
	return cookies
}

// Serves a file, with the name given by path, to the connection conn (although it could be any io.Writer). The info string contains extra info about how the serving went, for use in the log. E.g. "file not modified".
func ServeFile(conn io.Writer, path string, req *http.Request) (info string, err os.Error) {
	//fmt.Println("Attempting to serve file", req.URL.Path)
	// This is the actual code to serve the file. Should use http.ServeFile when we get this working, as it provides caching (304 not modified) and all kinds of other nifty features.
	file, e := os.Open(path)
	if e != nil {
		fmt.Println("Error opening file:", e)
		fmt.Fprintln(conn, "HTTP/1.1 404", e)
		fmt.Fprintln(conn, "Content-Length:", len(e.String()))
		fmt.Fprintln(conn, "")
		fmt.Fprintln(conn, e)
		return
	}
	finfo, e := file.Stat()
	if e != nil {
		fmt.Println("Error stating file:", e)
		fmt.Fprintln(conn, "HTTP/1.1 500", e)
		fmt.Fprintln(conn, "Content-Length: 0")
		fmt.Fprintln(conn, "\n")
		return
	}
	// Stolen from HTTP library
	if t, _ := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); t != nil && finfo.Mtime_ns/1e9 <= t.Seconds() {
		fmt.Fprintln(conn, "HTTP/1.1 304 Not Modified")
		fmt.Fprintln(conn, "")
		//fmt.Println("[File not modified]")
		info = "file not modified"
		return
	}
	ext := pathm.Ext(path)
	ctype := mime.TypeByExtension(ext)
	// No extension, serve as a generic binary
	// TODO: Add check for plain text
	if ctype == "" {
		ctype = "application/octlet-stream"
	}
	fmt.Println(ext, ctype)
	fmt.Println("Serving file", path)
	fmt.Fprintln(conn, "HTTP/1.1 200 OK")
	fmt.Fprintln(conn, "Last-Modified:", time.SecondsToUTC(finfo.Mtime_ns/1e9).Format(http.TimeFormat))
	fmt.Fprintln(conn, "Content-Type:", ctype)
	fmt.Fprintln(conn, "Content-Length:", finfo.Size)
	fmt.Fprintln(conn, "")
	io.Copy(conn, file)
	//http.ServeFile(c, r, r.URL.Path[1:])
	return
}
