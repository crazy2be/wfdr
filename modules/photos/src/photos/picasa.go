package main

// A shadow of it's former self. The actual "fetcher" is in ./picasa-updater.
// Loads JSON data structres from files in data/picasa, as well as handling user authentication and file upload requests (although the front-end is in the photo package).

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"time"
	// Local imports
	"github.com/crazy2be/iomod"
	"github.com/crazy2be/jsonutil"
	"github.com/crazy2be/user"
	"util/dlog"
	"util/picasa"
	"util/template"
)

// Defined like this so that moving the actual definitions to an external file doesn't break things.
type PhotoGroup picasa.PhotoGroup
type Photo picasa.Photo
type Album picasa.Album

// Cache the picasa JSON data structure.
var cachedAlbums []Album
var cachedPhotos = make(map[string][]Photo, 10)

// Refresh the cache every hour or so.
func init() {
	go func() {
		for {

			albumsFilename := "data/picasa/albums.json"
			jsonutil.DecodeFromFile(albumsFilename, &cachedAlbums)
			for _, album := range cachedAlbums {
				fmt.Println("album Link: ", album.Link, ", album ID: ", album.AlbumId)
				photosFilename := "data/picasa/albums/" + album.Link + ".json"
				var photos []Photo
				jsonutil.DecodeFromFile(photosFilename, &photos)
				cachedPhotos[album.Link] = photos
			}

			// Note, when setting this, that updating the feed
			// Causes lots of things to be swapped out in
			// Low-memory enviroments. Choose a balance.
			time.Sleep(1000e9) // 1000 seconds.

			// Why, you may ask, is this down here? Well, because i don't want to have to wait for all of the photos to load each time i start the module for debugging.
			updater := exec.Command("bin/picasa-updater")
			err := updater.Run()
			if err != nil {
				fmt.Println("Error running photos updater:", err)
				continue
			}
			updater.Wait()
		}
	}()
}

func GetAlbums() *[]Album {
	return &cachedAlbums
}

func GetPhotos(albumname string) []Photo {
	return cachedPhotos[albumname]
}

// For login authentication from picasa.
// TODO: Add error handling.
func AuthHandler(c http.ResponseWriter, r *http.Request) {
	// Get the token supplied in the URL.
	picasaLen := len("token=")
	url, _ := url.QueryUnescape(r.URL.RawQuery)
	token := url[picasaLen:]
	fmt.Println(token, r.URL.RawQuery)

	// Try to upgrade the token to a multi-use one. See
	// http://code.google.com/apis/accounts/docs/AuthSub.html
	req := picasa.NewRequest("https://www.google.com/accounts/accounts/AuthSubSessionToken", token, "GET")
	resp, e := picasa.Send(req)

	// Get the upgraded token value
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		fmt.Println(e)
	}
	resp.Body.Close()
	if len(body) <= picasaLen {
		dlog.Println("Invalid or missing token! Response received was:", body)
		template.Error500(c, r, nil)
	}
	upgradedToken := body[picasaLen:]
	fmt.Println("Upgraded Token: ", string(upgradedToken))

	// Finally, save the upgraded token in the server-side session.
	u, _ := user.Get(c, r)
	u.Set("picasa-authsub-token", string(upgradedToken))
	http.Redirect(c, r, "/photos/upload", http.StatusFound)
}

// Handles requests to upload to picasa in multipart/form-data format. Normally bound to /picasa/upload.
func UploadHandler(c http.ResponseWriter, r *http.Request) {
	// Handles multipart/form-data requests
	albumName, contentType, fileName, contentLength, fileReader, e := multipartUploadHandler(r)

	if e != nil {
		fmt.Println("Not multipart")

		// Handle a normal POST request, likely from the html5 uploader.
		albumName = r.Header.Get("X-Album-Name")
		contentType = r.Header.Get("Content-Type")
		fileName = r.Header.Get("Slug")
		contentLength = r.Header.Get("Content-Length")
		fileReader = r.Body
		defer r.Body.Close()
	}

	//var s *session.Session
	fmt.Println("Handling upload request at /picasa/upload.")
	fmt.Println(r)
	s, e := user.GetExisting(r)
	if e != nil {
		fmt.Fprintln(c, "Invaid session. Please login.")
		return
	}
	resp, e := uploadToPicasa(albumName, contentType, fileName, contentLength, s, fileReader)
	if e != nil {
		fmt.Println("Error uploading to picasa!", e)
		fmt.Fprintln(c, "Error uploading to picasa:", e)
	}
	handleError(resp, c)
	// For debugging.
	resp.Write(os.Stdout)
	return

	//TODO: Add Multipart/Form-Data support.
}

func multipartUploadHandler(r *http.Request) (albumName, contentType, fileName, contentLength string, fileReader io.Reader, err error) {
	mbound, err := checkMultipart(r)
	if err != nil {
		return
	}
	// Count reader, counts bytes read as they are read.
	creader := iomod.NewReadCounter(r.Body)
	mreader := multipart.NewReader(creader, mbound)

	sconlen := r.Header.Get("Content-Length")
	conlen, err := strconv.Atoi(sconlen)
	// Picasa REQUIRES Content-Length!
	if err != nil {
		fmt.Println("No Content-Length header or invalid value!", sconlen)
		return
	}

	for {
		var mpart *multipart.Part
		mpart, err = mreader.NextPart()
		if mpart != nil {
			fmt.Println("Multipart handler:", mpart, mpart.FormName(), err)
		} else {
			return
		}
		conlen -= 1
		name := mpart.FormName()
		switch name {
		case "album":
			var albumNameBytes []byte
			albumNameBytes, err = ioutil.ReadAll(mpart)
			if err != nil {
				fmt.Println("Error reading album name!", albumName, err)
				return
			}
			fmt.Println("Read", creader.Count, "bytes so far ( content-length is", r.Header["Content-Length"], ")")
			albumName = string(albumNameBytes)
		case "Filedata":
			contentType = mpart.Header.Get("Content-Type")

			var mtypes map[string]string
			_, mtypes, err = mime.ParseMediaType(mpart.Header.Get("Content-Disposition"))
			if err != nil {
				return
			}
			fileName = mtypes["filename"]

			fmt.Println("Read", creader.Count, "bytes so far ( content-length is", r.Header.Get("Content-Length"), ")")

			// We have to do this, because it seems like the only reliable way to determine the size of the file... Hopefully the files they send are not too large...
			// WARNING: Security vunerability with large files, could overrun the server.
			buf := new(bytes.Buffer)
			io.Copy(buf, mpart)
			fileReader = buf
			contentLength = strconv.Itoa(buf.Len())
		}
	}
	return
}

func checkMultipart(r *http.Request) (boundary string, err error) {
	v := r.Header.Get("Content-Type")
	if len(v) <= 0 {
		return "", errors.New("Not Multipart: Missing Content-Type header")
	}
	d, params, err := mime.ParseMediaType(v)
	if err != nil {
		return
	}
	if d != "multipart/form-data" {
		return "", errors.New("Not Multipart: MIME type is incorrect")
	}
	boundary, ok := params["boundary"]
	if !ok {
		return "", errors.New("Invalid Multipart: Missing boundary")
	}
	return
}

func uploadToPicasa(album, contentType, filename, size string, u *user.User, file io.Reader) (resp *http.Response, e error) {
	// Dial and send headers
	conn, e := picasa.Dial("picasaweb.google.com")
	// There is no sanitation of the user data here, but it should be fine since it's their own account anyway...
	fmt.Fprintln(conn, "POST http://picasaweb.google.com/data/feed/api/user/crazy1be/albumid/"+album+" HTTP/1.1")
	fmt.Fprintln(conn, "Authorization:", "AuthSub", "token=\""+u.Get("picasa-authsub-token")+"\"")
	fmt.Fprintln(conn, "Content-Type:", contentType)
	fmt.Fprintln(conn, "Content-length:", size)
	fmt.Fprintln(conn, "Slug:", filename)
	fmt.Fprintln(conn, "")
	// Copy the file data to the connection
	n, e := io.Copy(conn, file)
	fmt.Println(n, "bytes copied")
	fmt.Println("Sent request")
	// Get the response and close the connection
	resp, e = picasa.ReadFrom(conn, "POST")
	conn.Close()
	if e != nil {
		fmt.Println(e)
	}
	return
}

func handleError(resp *http.Response, c http.ResponseWriter) (e error) {
	//fmt.Printf("Response: %#v", resp)
	// 2xx is the "OK" http status range.
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Tell the uploader that upload was successfull.
		fmt.Fprintf(c, "1")
	} else {
		// TODO: Should return an acutal errror.
		// Respond with an error for the client side.
		fmt.Fprintf(c, "%d\n", resp.StatusCode)
		io.Copy(c, resp.Body)
	}
	return
}
