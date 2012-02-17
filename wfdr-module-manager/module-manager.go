// The go portion of the module management utility, used to start/stop/restart modules. This is simply a frontend for the wfdr-deamon, which actually manages the modules, starting and stopping them on demand.
package main

import (
	"flag"
	"fmt"
	"log"
	//"strings"
	// Local imports
	//"util/osutil"
	"util/moduled"
)

func main() {
	var action, name string

	flag.StringVar(&action, "action", "start", "What action would you like to take on this module? Valid options are start, stop, and restart.")
	flag.StringVar(&name, "modulename", "", "What module would you like to take this action on? Valid options are all, base, or the name of any module.")

	flag.Parse()

	if name == "" {
		log.Fatal("You must specify a module name.")
	}

	//if !osutil.FileExists("modules/" + name) {
	//	log.Fatal("Module not found, aborting.")
	//}

	switch action {
	case "start":
		StartModule(name)
	case "stop":
		StopModule(name)
	case "restart":
		StopModule(name)
		StartModule(name)
	default:
		log.Fatal("Unrecognized action ", action)
	}
}

func StartModule(name string) error {
	client, err := moduled.RPCConnect()
	if err != nil {
		log.Fatal(err.Error() + ". In all likelyhood, wfdr-deamon is not running. Start it and this error should go away! It's also possible that you don't have permission to talk to the deamon process, in which case i can't help you.")
	}
	var dummy int = 1000
	err = client.Call("ModuleSrv.Start", &name, &dummy)
	if err != nil {
		handleError(err, "starting", name)
	}
	return err
}

func StopModule(name string) error {
	client, err := moduled.RPCConnect()
	if err != nil {
		log.Fatal(err)
	}
	var dummy int = 1000
	err = client.Call("ModuleSrv.Stop", &name, &dummy)
	if err != nil {
		handleError(err, "stopping", name)
	}
	return err
}

func handleError(err error, verb, name string) {
// 	// WARNING: HACK HACK HACK
// 	if strings.Index(err.Error(), "Unmarshal") != -1 {
// 		return
// 	}
	fmt.Println("Error "+verb+" module", name, ":", err)

}
