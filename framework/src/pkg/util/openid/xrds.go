package openid

import (
	"xml"
	"fmt"
	"io"
	"strings"
)

type XRDSIdentifier struct {
        XMLName xml.Name "Service"
        Type []string
        URI string
        LocalID string
}
type XRD struct {
	XMLName xml.Name "XRD"
	Service XRDSIdentifier
}
type XRDS struct {
	XMLName xml.Name "XRDS"
	XRD XRD
}

func (o *OpenID) ParseXRDS(r io.Reader){
	XRDS := new(XRDS)
	err := xml.Unmarshal(r, XRDS)
	if err != nil {
		fmt.Printf(err.String())
                return
	}
	XRDSI := XRDS.XRD.Service

	XRDSI.URI = strings.TrimSpace(XRDSI.URI)
	XRDSI.LocalID = strings.TrimSpace(XRDSI.LocalID)

	fmt.Printf("%v\n", XRDSI)

	if StringTableContains(XRDSI.Type,"http://specs.openid.net/auth/2.0/server") {
		o.OPEndPoint = XRDSI.URI
		fmt.Printf("OP Identifier Element found\n")
	} else if StringTableContains(XRDSI.Type, "http://specs.openid.net/auth/2.0/signon") {
		fmt.Printf("Claimed Identifier Element found\n")
		o.OPEndPoint = XRDSI.URI
		o.ClaimedIdentifier = XRDSI.LocalID
	}
}


func StringTableContains (t []string, s string) bool {
	for _,v := range t {
		if v == s {
			return true
		}
	}
	return false
}
