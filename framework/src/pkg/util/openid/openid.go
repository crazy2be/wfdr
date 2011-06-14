package openid

import (
	"fmt"
	"io"
	"regexp"
	"os"
	"bytes"
	"http"
	"strings"
	)

const (
	_ = iota
	IdentifierXRI
	IdentifierURL
)

type OpenID struct {
	Identifier string  // Discovery
	IdentifierType int // Discovery
	OPEndPoint string  // Discovery
	ClaimedIdentifier string // Discovery
	OPLocalIdentifier string // Discovery
	Params map[string] string
	RPUrl string
	Hostname string
	Request string
	Realm string
	ReturnTo string
}

func (o *OpenID) normalization() {
	//1.  If the user's input starts with the "xri://" prefix, it MUST be stripped off, so that XRIs are used in the canonical form.
	if strings.HasPrefix(o.Identifier, "xri://") {
		o.Identifier = o.Identifier[6:]
	}

	// 2. If the first character of the resulting string is an XRI Global Context Symbol ("=", "@", "+", "$", "!") or "(", as defined in Section 2.2.1 of [XRI_Syntax_2.0] (Reed, D. and D. McAlpin, “Extensible Resource Identifier (XRI) Syntax V2.0,” .), then the input SHOULD be treated as an XRI.
	if o.Identifier[0] == '=' || o.Identifier[0] == '@' || o.Identifier[0] == '+' || o.Identifier[0] == '$' || o.Identifier[0] == '!'  {
		o.IdentifierType = IdentifierXRI
		fmt.Printf("It is an XRI\n")
		return
	}

	// 3. Otherwise, the input SHOULD be treated as an http URL; if it does not include a "http" or "https" scheme, the Identifier MUST be prefixed with the string "http://". If the URL contains a fragment part, it MUST be stripped off together with the fragment delimiter character "#". See Section 11.5.2 (HTTP and HTTPS URL Identifiers) for more information.
	o.IdentifierType = IdentifierURL
	if ! strings.HasPrefix(o.Identifier, "http://") && ! strings.HasPrefix(o.Identifier, "https://") {
		o.Identifier = "http://" + o.Identifier
	}

	// 4. URL Identifiers MUST then be further normalized by both following redirects when retrieving their content and finally applying the rules in Section 6 of [RFC3986] (Berners-Lee, T., “Uniform Resource Identifiers (URI): Generic Syntax,” .) to the final destination URL. This final URL MUST be noted by the Relying Party as the Claimed Identifier and be used when requesting authentication (Requesting Authentication).
}

func (o *OpenID) discovery() os.Error {
	//1.  If the identifier is an XRI, [XRI_Resolution_2.0]  (Wachob, G., Reed, D., Chasen, L., Tan, W., and S. Churchill, “Extensible Resource Identifier (XRI) Resolution V2.0 - Committee Draft 02,” .)  will yield an XRDS document that contains the necessary information. It should also be noted that Relying Parties can take advantage of XRI Proxy Resolvers, such as the one provided by XDI.org at http://www.xri.net. This will remove the need for the RPs to perform XRI Resolution locally.
	if o.IdentifierType == IdentifierXRI {
		fmt.Printf("XRI Discovery not implemented yet\n")
	}

	//2. If it is a URL, the Yadis protocol (Miller, J., “Yadis Specification 1.0,” .) [Yadis] SHALL be first attempted. If it succeeds, the result is again an XRDS document.
	if o.IdentifierType == IdentifierURL {
		r,err := Yadis(o.Identifier)
		if err != nil {
			return err
		}
		o.ParseXRDS(r)
	}

	//3. If the Yadis protocol fails and no valid XRDS document is retrieved, or no Service Elements (OpenID Service Elements) are found in the XRDS document, the URL is retrieved and HTML-Based discovery (HTML-Based Discovery) SHALL be attempted.
	// Not Yet implemented



	// If the end user entered an OP Identifier, there is no Claimed Identifier. For the purposes of making OpenID Authentication requests, the value "http://specs.openid.net/auth/2.0/identifier_select" MUST be used as both the Claimed Identifier and the OP-Local Identifier when an OP Identifier is entered.
	if o.ClaimedIdentifier == "" {
		fmt.Printf("Set identifier_select\n")
		o.OPLocalIdentifier = "http://specs.openid.net/auth/2.0/identifier_select"
		o.ClaimedIdentifier = "http://specs.openid.net/auth/2.0/identifier_select"
	}

	return nil
}

func mapToUrlEnc (params map[string] string) string {
	url := ""
	for k,v := range (params) {
		url = fmt.Sprintf("%s&%s=%s",url,k,v)
	}
	return url[1:]
}

func urlEncToMap (url string) map[string] string {
	// We don't know how elements are in the URL so we create a list first and push elements on it
	pmap := make(map[string] string)
	//url,_ = http.URLUnescape(url)
	var start, end, eq, length int
	length = len(url)
	start = 0
	for start < length && url[start] != '?' { start ++ }
	end = start
	for end < length {
		start = end + 1
		eq = start
		for eq < length && url[eq] != '=' { eq++ }
		end = eq + 1
		for end < length && url[end] != '&' { end++ }
	
		fmt.Printf("Trouve: %s : %s\n", url[start:eq], url[eq+1:end])
		pmap[url[start:eq]] = url[eq+1:end]
	}
	return pmap
}

func (o *OpenID) GetUrl() (string, os.Error) {
	o.normalization()
	err := o.discovery()
	if err != nil {	return "", err }


	params := map[string] string {
		"openid.ns": "http://specs.openid.net/auth/2.0",
		"openid.mode" : "checkid_setup",
	}
	if o.Realm != "" {
		params["openid.realm"] = o.Realm
	}
	if o.Realm != "" && o.ReturnTo != "" {
		params["openid.return_to"] = o.Realm + o.ReturnTo
	}
	if o.ClaimedIdentifier != "" {
		params["openid.claimed_id"] = o.ClaimedIdentifier
		if o.OPLocalIdentifier != "" {
			params["openid.identity"] = o.OPLocalIdentifier
		} else {
			params["openid.identity"] = o.ClaimedIdentifier
		}
	}
	return o.OPEndPoint + "?" + mapToUrlEnc(params), nil
}

func (o *OpenID) Verify() (grant bool, err os.Error) {
	grant = false
	err = nil
	

	// The value of "openid.return_to" matches the URL of the current request
	// if ! MExists(o.Params, "openid.return_to") {
	// 	err = os.ErrorString("The value of 'openid.return_to' is not defined")
	// 	return
	// }
	// if (fmt.Sprintf("%s%s", o.Hostname, o.Request) != o.Params["openid.return_to"]) {
	// 	err = os.ErrorString("The value of 'openid.return_to' does not match the URL of the current request")
	// 	return
	// }

	// Discovered information matches the information in the assertion

	// An assertion has not yet been accepted from this OP with the same value for "openid.response_nonce"

	// The signature on the assertion is valid and all fields that are required to be signed are signed
	grant, err = o.VerifyDirect()
	
	return
}

func (o *OpenID) ParseRPUrl(url string) {
	o.Params = urlEncToMap(url)
}

func (o *OpenID) VerifyDirect() (grant bool, err os.Error) {
	grant = false
	err = nil

	o.Params["openid.mode"] = "check_authentication"

	headers := make(http.Header, 1)
	headers.Add("Content-Type", "application/x-www-form-urlencoded")
	
	url,_ := http.URLUnescape(o.Params["openid.op_endpoint"])
	fmt.Printf("Verification: %s\nParams: %s\n",url, mapToUrlEnc(o.Params))
	r,error := post(url,
		headers,
		bytes.NewBuffer([]byte(mapToUrlEnc(o.Params))))
	if error != nil {
		fmt.Printf("erreur: %s\n", error.String())
		err = error
		return
	}
	fmt.Printf("Post done\n")
	if (r != nil) {
		buffer := make([]byte, 1024)
		fmt.Printf("Buffer created\n")
		io.ReadFull(r.Body, buffer)
		fmt.Printf("Body extracted: %s\n", buffer)
		grant, err = regexp.Match("is_valid:true", buffer)
		fmt.Printf("Response: %v\n", grant)
	}else {
		err = os.ErrorString("No response from POST verification")
		return
	}

	return
}
