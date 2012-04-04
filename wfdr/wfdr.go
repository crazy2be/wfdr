// A simple program that provides aliases for other programs in the wfdr toolchain.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"wfdr/moduled"

// 	"github.com/kylelemons/goat"
)

func main() {
	// 	buf := make([]byte, 3)
	state, err := moduled.SttyCbreak()
	defer state.Undo()
	if err != nil {
		log.Println(state)
		log.Fatal(err)
	}
	s := moduled.NewShell(os.Stdin, os.Stdout)
	for {
		cmd, err := s.ReadCommand()
		if err != nil {
			fmt.Println("Error while reading command: ", err)
		}
		if cmd != nil {
			fmt.Printf("Read command: %#v\n", cmd)
		}
	}
	os.Exit(0)
	if len(os.Args) < 2 {
		printHelp()
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
			moduleAction(fi.Name(), action)
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

func printHelp() {
	fmt.Println("Usage: wfdr <action> [<modulename>]...")
	fmt.Println(" Action can be one of stop, start, restart, compile, recompile, status, or list.")
	fmt.Println(" (status is not implemented)")
	fmt.Println(" modulename can be the name of any installed module, or 'all' to take an action on all modules. Multiple names can be specified, seperated by spaces. (e.g. wfdr start auth base main)")
}

func mustRun(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err == nil {
		return
	}
	ws, ok := err.(*exec.ExitError)
	if !ok {
		os.Exit(1)
	}
	if ws.Success() {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func moduleAction(module, action string) {
	var err error
	os.Setenv("PATH", os.Getenv("PATH")+":framework/sh:framework/bin")
	switch action {
	case "stop", "start", "restart":
		client, err := moduled.RPCConnect()
		if err != nil {
			log.Fatal(err)
		}
		if action == "stop" || action == "restart" {
			err = client.Stop(module)
		}
		if err != nil {
			break
		}
		if action == "start" || action == "restart" {
			err = client.Start(module)
		}
	case "compile":
		mustRun("wfdr-compile", module)
	case "recompile":
		mustRun("wfdr-compile", module, "-recompile")
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
			name := fi.Name()
			running, err := client.Status(name)
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
		printHelp()
		os.Exit(1)
	}
	if err != nil {
		log.Fatal(err)
	}
}
