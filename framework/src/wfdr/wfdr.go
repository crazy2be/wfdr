// A simple program that provides aliases for other programs in the wfdr toolchain.

package main


import (
	"os"
	"log"
	"fmt"
	"io/ioutil"
	// Local imports
	"util/osutil"
	"util/moduled"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wfdr <action> [<modulename>]...")
		fmt.Println("	Action can be one of stop, start, restart, compile, recompile, status, or list.")
		fmt.Println("	(status is not implemented)")
		fmt.Println("	modulename can be the name of any installed module, or 'all' to take an action on all modules. Multiple names can be specified, seperated by spaces. (e.g. wfdr start auth base main)")
		os.Exit(0)
	}
	action := os.Args[1]
	module := ""
	if len(os.Args) >= 3 {
		module = os.Args[2]
		//log.Fatal("Module name not provided...")
	}
	
	if module == "all" {
		fis, err := ioutil.ReadDir("modules")
		if err != nil {
			log.Fatal("Error opening modules directory, cannot possibly perform an action on 'all':", err)
		}
		for _, fi := range fis {
			moduleAction(fi.Name, action)
		}
		return
	}
	
	// Multiple module names specified
	if len(os.Args) > 3 {
		for i := 2; i < len(os.Args); i++ {
			moduleAction(os.Args[i], action)
		}
		return
	}
	
	moduleAction(module, action)
}

func moduleAction(module, action string) {
	var err os.Error
	os.Setenv("PATH", os.Getenv("PATH")+":framework/sh:framework/bin")
	switch action {
		case "stop", "start", "restart":
			_, err = osutil.WaitRun("wfdr-module-manager", []string{"-action="+action, "-modulename="+module})
		case "compile":
			_, err = osutil.WaitRun("wfdr-compile", []string{module})
		case "recompile":
			_, err = osutil.WaitRun("wfdr-compile", []string{module, "-recompile"})
		case "status":
			fmt.Println("Not implemented!")
		case "list":
			fis, err := ioutil.ReadDir("modules")
			if err != nil {
				log.Fatal("Error opening modules directory, cannot list modules.")
			}
			client, err := moduled.RPCConnect()
			if err != nil {
				log.Fatal(err)
			}
			for _, fi := range fis {
				name := fi.Name
				var running bool
				err = client.Call("ModuleSrv.Status", &name, &running)
				if err != nil {
					log.Fatal(err)
				}
				if running {
					fmt.Printf(" * %s\n", name)
				} else {
					fmt.Printf("   %s\n", name)
				}
			}
			err = client.Close()
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal("Unrecognized command ", action)
	}
	if err != nil {
		log.Fatal(err)
	}
}