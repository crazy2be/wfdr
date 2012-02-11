// A deamon that manages module processes, attempting to properly deal with them when they crash, and implementing manual starting and stopping of modules.
package main

import (
	"flag"
	"io"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"os/signal"
	"path"
	// Local imports
	"github.com/crazy2be/osutil"
	"util/dlog"
	"util/pipes"
)

func main() {
	// Kinda hackish :/
	wd, _ := os.Getwd()
	os.Setenv("PATH",
		path.Join(wd, "framework/bin")+":"+
			path.Join(wd, "framework/sh")+":"+
			os.Getenv("PATH"))

	os.MkdirAll("cache", 0755)

	var inpipe, outpipe, context string

	flag.StringVar(&inpipe, "inpipe", "cache/wfdr-deamon-pipe-in", "Name of a file that should be used for a input IPC pipe. Should not exist.")
	flag.StringVar(&outpipe, "outpipe", "cache/wfdr-deamon-pipe-out", "Name of a file that should be used for an output IPC pipe. Should not exist.")
	flag.StringVar(&context, "context", "debug", "Context to run the daemon and child processes from. Valid choices are 'debug', 'test' and 'prod'.")

	flag.Parse()

	if context != "debug" && context != "test" && context != "prod" {
		log.Println("Invalid context argument provided!")
		log.Fatal(flag.Lookup("context").Usage)
	}

	os.Setenv("WFDR_CONTEXT", context)

	if osutil.FileExists(inpipe) || osutil.FileExists(outpipe) {
		log.Fatal("Pipe files already exist, the daemon is likely already running. However, it is also possible that the daemon was not cleanly shut down on its last run, and the files linger. If you suspect this to be the case, remove cache/wfdr-deamon-pipe-in and cache/wfdr-deamon-pipe-out, then try starting the daemon again.")
	}

	infile, err := pipes.MakeAndOpen(inpipe)
	if err != nil {
		log.Fatal(err)
	}
	outfile, err := pipes.MakeAndOpen(outpipe)
	if err != nil {
		log.Fatal(err)
	}

	rwc := &pipes.PipeReadWriteCloser{infile, outfile}
	//dlog.Println("Made ReadWriteCloser")

	go monitorPipe(rwc)

	for {
		sig := <-signal.Incoming
		switch sig.(os.UnixSignal) {
		// SIGINT, SIGKILL, SIGTERM
		case 0x02, 0x09, 0xf:
			Exit(0)
		// SIGCHLD
		case 0x11:
		// Do nothing
		default:
			dlog.Println(sig)
			break
		}
	}
}

func Exit(status int) {
	for _, module := range modules {
		module.Stop()
	}
	os.Remove("cache/wfdr-deamon-pipe-in")
	os.Remove("cache/wfdr-deamon-pipe-out")
	os.Exit(status)
}

// Defined for the RPC package
type ModuleSrv int

// Starts the module. ret is never changed.
func (m *ModuleSrv) Start(name *string, ret *int) error {
	dlog.Println("Starting module:", *name)
	*ret = 12390
	_, err := StartModule(*name)
	return err
}

func (m *ModuleSrv) Stop(name *string, ret *int) error {
	dlog.Println("Stopping module:", *name)
	*ret = 12048
	err := StopModule(*name)
	return err
}

func (m *ModuleSrv) Status(name *string, running *bool) error {
	//var isrunning bool
	//*running = isrunning
	//dlog.Println(modules)
	mod, err := GetModule(*name)
	if err != nil {
		*running = false
		return nil
	}
	dlog.Println("Module was running last we checked...")
	*running = mod.IsRunning()
	if *running {
		dlog.Println("Module appears to be running!")
	}
	return nil
}

func monitorPipe(rwc io.ReadWriteCloser) {

	StartSharedSync()

	serv := rpc.NewServer()
	codec := jsonrpc.NewServerCodec(rwc)
	//dlog.Println("Made RPC server")
	m := new(ModuleSrv)
	serv.Register(m)
	//dlog.Println("Registered module service")
	dlog.Println("RPC server is started and awaiting commands!")
	serv.ServeCodec(codec)
	//dlog.Println("Serving on connection rwc")
}
