package openid

import (
	"errors"
	"fmt"
	"net/http"

	"encoding/xml"
	"io"
	"strings"
)

func searchHTMLMeta(r io.Reader) (string, error) {
	parser := xml.NewParser(r)
	var token xml.Token
	var err error
	for {
		token, err = parser.Token()
		if token == nil || err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		switch token.(type) {
		case xml.StartElement:
			if token.(xml.StartElement).Name.Local == "meta" {
				// Found a meta token. Verify that it is a X-XRDS-Location and return the content
				var content string
				var contentE bool
				var httpEquivOK bool
				contentE = false
				httpEquivOK = false
				for _, v := range token.(xml.StartElement).Attr {
					if v.Name.Local == "http-equiv" && v.Value == "X-XRDS-Location" {
						httpEquivOK = true
					}
					if v.Name.Local == "content" {
						content = v.Value
						contentE = true
					}
				}
				if contentE && httpEquivOK {
					return fmt.Sprint(content), nil
				}
			}
		}
	}
	return "", errors.New("Value not found")
}

func Yadis(url string) (io.Reader, error) {
	fmt.Printf("Search: %s\n", url)
	var headers http.Header
	headers.Add("Accept", "application/xrds+xml")

	r, err := get(url, headers)
	if err != nil || r == nil {
		fmt.Printf("Yadis: Error in GET\n")
		return nil, err
	}

	// If it is an XRDS document, parse it and return URI
	content := r.Header.Get("Content-Type")
	if content != "" && strings.HasPrefix(content, "application/xrds+xml") {
		fmt.Printf("Document XRDS found\n")
		return r.Body, nil
	}

	// If it is an HTML doc search for meta tags
	content = r.Header.Get("Content-Type")
	if content != "" && content == "text/html" {
		fmt.Printf("Document HTML found\n")
		url, err := searchHTMLMeta(r.Body)
		if err != nil {
			return nil, err
		}
		return Yadis(url)
	}

	// If the response contain an X-XRDS-Location header
	xrds := r.Header.Get("X-Xrds-Location")
	if xrds != "" {
		return Yadis(xrds)
	}

	// If nothing is found try to parse it as a XRDS doc
	return nil, nil
}
