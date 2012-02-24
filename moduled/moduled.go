// Provides functions for communicating with the module deamon through the rpc module.
package moduled

import (
	"os"
	"log"
)


var knownLayouts []string = []string{"mobile", "desktop"}

type CacheMonitor struct {
	source *os.File
	dest   *os.File
	dlog   *log.Logger
	typ    string
}

func NewCacheMonitor(source string, dest string, typ string) (*CacheMonitor, error) {
	var err error
	cm := new(CacheMonitor)
	
	cm.source, err = os.Open(source)
	if err != nil {
		return nil, err
	}
	
	err = os.MkdirAll(dest, 0744)
	if err != nil {
		return nil, err
	}
	
	cm.dest, err = os.Open(dest)
	if err != nil {
		return nil, err
	}
	
	switch typ {
		case "css", "js", "img", "tmpl":
		default:
			log.Println("wfdr/moduled: Warning: Unrecognized file type", typ)
	}
	cm.typ = typ
	
	devnull, err := os.Open(os.DevNull)
	if err != nil {
		return nil, err
	}
	cm.dlog = log.New(devnull, "wfdr/moduled/CacheMonitor:", log.LstdFlags)
	
	return cm, nil
}

