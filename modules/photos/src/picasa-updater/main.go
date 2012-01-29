// Updates the photo lists from picasa, translating the picasa-style JSON to include only the fields that the photos module requires.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	// Local imports
	"github.com/crazy2be/iomod"
	"github.com/crazy2be/jsonutil"
	"util/picasa"
)

// The URL to request the albums from. Note that ?alt=json is added to this.
const userFeedURL = "http://picasaweb.google.com/data/feed/api/user/default"

// The time format for the updated and published fields from the json feed.
const jsonTimeFormat = "2006-01-02T15:04:05:00.000Z"

// The authentication token for picasa, shouldn't really be hard coded. However, it works for now, and i can't really think of a much better way.
const picasaAuthSubToken = ""

// The data structure for the JSON feed from picasa. Add and remove fields as
// neccessary, We should avoid having too many unnecesary fields for the sake
// of confusing ourselves! Note that $ is replaced with S compared to a normal
// picasa feed, because go doesn't like $ in variable names.
type JSONPhoto struct {
	Title struct {
		ST   string
		Type string
	}
	Summary struct {
		ST   string
		Type string
	}
	RelLink struct {
		Root     string // TODO
		Relative string // relitive to photos page
	}
	Link []struct {
		Rel  string
		Type string
		Href string
	}
	GphotoSname struct {
		ST string
	}
	GphotoSid struct {
		ST string
	}
	GphotoStimestamp struct {
		ST string
	}
	Updated struct {
		ST string
	}
	MediaSgroup struct {
		MediaScontent []struct {
			Url    string
			Type   string
			Medium string
			Width  int
			Height int
		}
		MediaScredit []struct {
			ST string
		}
		MediaSdescription struct {
			ST   string
			Type string
		}
		MediaSthumbnail []struct {
			Url    string
			Height int
			Width  int
		}
	}
}
type UserFeed struct {
	Feed struct {
		Entry []JSONPhoto
	}
}
type AlbumFeed struct {
	Feed struct {
		Entry []JSONPhoto
	}
}

// Defined like this so that moving the actual definitions to an external file doesn't break things.
type PhotoGroup picasa.PhotoGroup
type Photo picasa.Photo
type Album picasa.Album

// Grab the photo info
func main() {
	// Makes the data directory first if it does not exist
	os.MkdirAll("data/picasa/albums", 0755)
	os.MkdirAll("data/picasa/albumsHD", 0755)

	var albumName string
	flag.StringVar(&albumName, "album", "all", "Specify a specific album to update, according to picasa ID. This is a series of numbers that probably looks something like 5510949320. You can optionally also specify a name, such as \"WinterFormalCouples2010\" If you only want to update the index, specify \"list\" in place of an ID. If you want to update all albums, you don't have to specify this command-line option.")
	flag.Parse()

	if albumName == "all" {
		fmt.Println("::Updating All Albums::")
		albums, albumsHD := updateList()
		for _, album := range albums {
			updateAlbum(album.AlbumId, album.Link)
		}
		for _, album := range albumsHD {
			updateAlbum(album.AlbumId, album.Link)
		}
		return
	}
	if albumName == "list" {
		fmt.Println("::Updating Album List::")
		updateList()
		return
	}
	fmt.Printf("::Attempting to Update Album %s::\n", albumName)

	var albums []picasa.Album
	albumsFilename := "data/picasa/albums.json"
	jsonutil.DecodeFromFile(albumsFilename, &albums)

	for _, album := range albums {
		if album.AlbumId == albumName ||
			album.Link == albumName {
			updateAlbum(album.AlbumId, album.Link)
		}
	}
}

func updateList() (albums []Album, albumsHD []Album) {
	var err error
	albums, albumsHD, err = getAlbums()
	if err != nil {
		log.Fatal("Error updating album list: ", err)
	}
	albumsFilename := "data/picasa/albums.json"
	albumsHDFilename := "data/picasa/albumsHD.json"
	jsonutil.EncodeToFile(albumsFilename, albums)
	jsonutil.EncodeToFile(albumsHDFilename, albumsHD)
	return
}

func updateAlbum(id, link string) {
	fmt.Println("Updating album", link, "( id", id, ")")
	photos, _ := getPhotos(id)
	photosFilename := ""
	if isHD(link) {
		photosFilename = "data/picasa/albumsHD/" + link + ".json"
	} else {
		photosFilename = "data/picasa/albums/" + link + ".json"
	}
	jsonutil.EncodeToFile(photosFilename, photos)
}

func getAlbums() (albums []Album, albumsHD []Album, e error) {
	feed, e := getUserFeed()
	albumNum := len(feed.Feed.Entry)
	albumNumHD := 0
	for _, entry := range feed.Feed.Entry {
		if isHD(entry.GphotoSname.ST) {
			albumNum--
			albumNumHD++
		}
	}
	albums = make([]Album, albumNum)
	albumsHD = make([]Album, albumNumHD)
	if e != nil {
		fmt.Printf("Error: %#v\n", e)
		return
	}
	fmt.Printf("Number of albums: %d\n", len(albums))
	// Translate the values in the JSON structure to the saner values
	// in the struct for external usage (Album).
	i := 0
	j := 0
	for _, entry := range feed.Feed.Entry {
		var album Album //= &albums[i]
		album.Title = entry.Title.ST
		album.Summary = entry.Summary.ST
		album.AlbumId = entry.GphotoSid.ST
		album.Link = entry.GphotoSname.ST
		album.Url = entry.MediaSgroup.MediaScontent[0].Url
		//var modified, _ = parsetime.Parse(entry.Updated.ST)
		//fmt.Printf("Modified: %#v\n", modified)
		if isHD(album.Link) {
			albumsHD[j] = album
			j++
		} else {
			albums[i] = album
			i++
		}
	}
	return
}

func isHD(albumLink string) bool {
	return albumLink[len(albumLink)-2:] == "HD"
}

func getPhotos(albumid string) (photos []Photo, e error) {
	feed, e := getAlbumFeed(albumid)
	photos = make([]Photo, len(feed.Feed.Entry))
	if e != nil {
		fmt.Println("Error:", e)
		return
	}
	fmt.Printf("Number of photos in album %s: %d\n", albumid, len(photos))
	hasCover := false
	for i, entry := range feed.Feed.Entry {
		if strings.ToLower(entry.Title.ST) == "cover.jpg" {
			hasCover = true
			continue
		}
		var photo = &photos[i]
		photo.Title = entry.Title.ST
		photo.Summary = entry.Summary.ST
		var g = entry.MediaSgroup
		var p = g.MediaScontent[0]
		photo.Url = p.Url
		photo.Width = p.Width
		photo.Height = p.Height
		photo.Thumbnails.Small = g.MediaSthumbnail[0]
		photo.Thumbnails.Medium = g.MediaSthumbnail[1]
		photo.Thumbnails.Large = g.MediaSthumbnail[2]
		var modified, e = picasa.ParseTimestamp(entry.Updated.ST)
		if e != nil {
			fmt.Printf("WARNING: Error parsing time string %s with format %s: %s", entry.Updated.ST, jsonTimeFormat, e)
		} else {
			photo.Modified = modified
		}
		timestamp, e := strconv.ParseInt(entry.GphotoStimestamp.ST, 10, 64)
		photo.Published = time.Unix(timestamp/1000, 0)

		//fmt.Printf("Published: %s\n", photo.Published)
		if e != nil {
			fmt.Printf("WARNING: Error parsing time string %s into local date: %s", entry.GphotoStimestamp.ST, e)
		}
	}
	if hasCover {
		photos = photos[:len(photos)-1]
	}
	return
}

// This one is internal.
func getUserFeed() (UserFeed, error) {
	// Replace $ with S in the incoming JSON... Go doesn't like $ in variable names
	var rep = iomod.NewReplacer(getUserFeedReader(), "\\$", "S")
	var dec = json.NewDecoder(rep)
	var feed UserFeed
	var e = dec.Decode(&feed)
	if e != nil {
		fmt.Println("Getting user feed")
		fmt.Println("Error:", e)
	}
	return feed, e
}

func getAlbumFeed(name string) (feed AlbumFeed, e error) {
	var rawFeed = iomod.NewReplacer(getAlbumFeedReader(name), "\\$", "S")
	var feedDecoder = json.NewDecoder(rawFeed)
	e = feedDecoder.Decode(&feed)
	if e != nil {
		fmt.Println("Retrieving album feed for", name)
		fmt.Println("Error:", e)
	}
	return
}

func getUserFeedReader() io.Reader {
	return getFeedReader("")
}

func getFeedReader(name string) io.Reader {
	// Dial and send headers
	conn, e := picasa.Dial("picasaweb.google.com")

	Path := ""
	QueryString := "alt=json"
	// We are requesting an album, not the list of albums.
	if name != "" {
		Path = "/albumid/" + name
	} else {
		QueryString = "type=album&" + QueryString
	}
	URL := userFeedURL + Path + "?" + QueryString

	fmt.Fprintln(conn, "GET "+URL+" HTTP/1.1")
	fmt.Fprintln(conn, "Authorization:", "AuthSub", "token=\""+picasaAuthSubToken+"\"")
	fmt.Fprintln(conn, "")

	//fmt.Println("Sent request")
	// Get the response. We can't close the connection yet, or the data will not all be read.
	resp, e := picasa.ReadFrom(conn, "GET")
	if e != nil {
		fmt.Println("Error!", e)
	}

	fmt.Printf("Requesting URL %s from picasa\n", URL)

	return resp.Body
}

func getAlbumFeedReader(name string) io.Reader {
	return getFeedReader(name)
}
