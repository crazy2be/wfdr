/* A simple program that is run at the initialization of the server, and runs until the server stops. It looks at all the files in css/ and js/, and then minifies and does preprocessing on them if they have been modified since their corresponding cache files were last updated.

It remains running in order to capture inotify events, and then automatically update the cache folder when css/ files are changed. This is designed to ease development, so that a full server restart is not required on each file change.

Please note that one instance is run for each folder that needs monitoring (e.g. one for modules/events/css, one for modules/events/js, one for modules/events/tmpl, one for modules/pages/css...)
*/
// TODO: Hook into the inotify osx equivalent.
package main


import (
	"os"
	"fmt"
	"log"
	"flag"
	"strings"
	"os/signal"
	"os/inotify"
	"path/filepath"
	// Local imports
	"github.com/crazy2be/osutil"
)

var dlog *log.Logger

func main() {
	sourceDir, destDir, destSubfolder, fileType, mode := "", "", "", "", ""

	var debug bool

	flag.StringVar(&sourceDir, "sourcedir", "", "The source directory where this cache-monitor should pull files from. The cache manager automatically watches files in base/, mobile/, and desktop/, calling the handler for each modified file.")
	flag.StringVar(&destDir, "destdir", "", "The cache directory where this cache-monitor should place files in. If this directory does not exist, it is created. Subfolders for each different layout are generated automatically.")
	flag.StringVar(&destSubfolder, "destsubfolder", "", "Do you want the files to end up in a subdirectory? This is different from changing the destdir in that the subfolder is applied to the path after the layout. E.g. setting the subfolder to 'events' would make the files end up in $sourcedir/$layout/$events/$file for each layout.")
	flag.StringVar(&fileType, "filetype", "", "What filetype is this cache-monitor monitoring? This generally, but not always, corresponds to the file extension. For example, 'js', 'css', 'img', and 'tmpl' are all valid types.")
	flag.StringVar(&mode, "mode", "sync", "What mode should the monitor run in? Valid options include:\n\tsync: Runs once, performing a one-way sync of all files.\n\tdeamon: Runs in the background, updating files as they are changed. Currently only supported on Linux.")
	flag.BoolVar(&debug, "debug", false, "Run in debug mode? Outputs a lot more garbage.")
	flag.Parse()
	

	if debug {
		dlog = log.New(os.Stderr, "DEBUG: ", log.Ltime|log.Lshortfile)
	} else {
		nul, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
		if err != nil {
			log.Fatal("Could not open null device? Wtf?")
		}
		dlog = log.New(nul, "", 0)
	}

	if mode == "deamon" {
		dlog.Println("Warning: Deamon mode is still beta. Currently it only works on linux systems with inotify support.")
	}

	var v Visitor
	v.source = new(Path)
	v.dest = new(Path)
	v.source.root = sourceDir
	v.dest.root = destDir
	v.dest.subfolder = destSubfolder
	v.fileType = fileType
	
	var err os.Error
	
	if mode == "deamon" {
		v.deamon = true
		v.watcher, err = inotify.NewWatcher()
		if err != nil {
			log.Fatal("Error creating inotify instance; deamon mode cannot possibly work.")
		}
	}

	errchan := make(chan os.Error, 100)
	donechan := make(chan bool, 1)
	
	go func() {
		for {
			select {
				case err := <-errchan:
					fmt.Println("Error in walker:", err)
				case <-donechan:
					return
			}
		}
	}()
	
	filepath.Walk(sourceDir, &v, errchan)
	donechan <-true
	
	if v.deamon {
		v.InotifyLoop()
	}
}

func (v *Visitor) String() string {
	return fmt.Sprintf("%#v", v)
}

type Visitor struct {
	source, dest *Path // Path of the current file
	fileType     string // Filetype of the current file.
	deamon       bool // Is the visitor operating in deamon mode? If true, will register hooks rather than check last-modified dates.
	watcher      *inotify.Watcher // Watcher, nil if deamon is false.
}

func (v *Visitor) InotifyLoop() {
	//log.Println("Watching for inotify events...")
	//log.Println("What")
	for {
		select {
			case ev := <-v.watcher.Event:
				//os.Exit(0)
				v.WatcherEvent(ev)
			case sig := <-signal.Incoming:
				switch (sig.(os.UnixSignal)) {
				// SIGINT, SIGKILL, SIGTERM
				case 0x02, 0x09, 0xf:
					v.watcher.Close()
					os.Exit(0)
				// SIGCHLD
				case 0x11:
					// Do nothing
				default:
					log.Println("Unhandled signal:", sig)
				}
			case err := <-v.watcher.Error:
				log.Println("Error in watcher:", err)
		}
	}
}

func (v *Visitor) WatcherEvent(ev *inotify.Event) {
	// Ignore qt temp files (from kate)
	if strings.Index(strings.ToLower(ev.Name), "qt_temp") != -1 {
		dlog.Println("Ignoring kate temp file...")
		return
	}
	// Ignore temp/backup files.
	if strings.HasSuffix(ev.Name, "~") {
		dlog.Println("Ignoring backup file (suffix ~)")
		return
	}
	dlog.Println("Got watcher event for file:", ev.Name)
	v.source = CalcPath(v.source.root, v.source.subfolder, ev.Name)
	v.ReloadCurrent()
	//log.Println(ev)
}

func (v *Visitor) VisitDir(dirpath string, fi *os.FileInfo) bool {
	// Only do something if in deamon mode.
	if !v.deamon {
		return true
	}
	err := v.watcher.AddWatch(dirpath, inotify.IN_MODIFY|inotify.IN_CREATE|inotify.IN_DELETE|inotify.IN_MOVE)
	if err != nil {
		dlog.Println("WARNING: Watcher failed to add a watch for directory", dirpath, ":", err, "(you will have to manually reload the module for your changes to take effect)")
		return true
	}
	dlog.Println("Watching directory", dirpath, "for changes")
	return true
}

func (v *Visitor) VisitFile(fpath string, fi *os.FileInfo) {
	// Don't do anything for files when in deamon mode.
	if v.deamon {
		return
	}
	dlog.Println("Checking file", fpath)
	v.source = CalcPath(v.source.root, v.source.subfolder, fpath)
	if v.source.layout == "base" {
		v.CheckFile("desktop", fi)
		v.CheckFile("mobile", fi)
		return
	}
	v.CheckFile(v.source.layout, fi)
	return
}

func (v *Visitor) CheckFile(layout string, fi *os.FileInfo) {
	// Gross
	fpath := v.source.JoinWithRootLayoutAndSubfolder(v.dest.root, layout, v.dest.subfolder)
	cachefi, err := os.Stat(fpath)
	if err != nil {
		// we always assume the error is that the cache file doesn't exist. Perhaps this is not wise?
		//fmt.Println("Error stating cache file:", err)
		v.ReloadFile(layout)
		return
	}
	if cachefi.Mtime_ns < fi.Mtime_ns {
		fmt.Println("Reloading file:", fpath)
		v.ReloadFile(layout)
	} else {
		dlog.Println("Shouldn't reload file:", fpath)
	}
}

func (v *Visitor) ReloadCurrent() {
	dlog.Println("Reloading file:", v.source)
	if v.source.layout == "base" {
		v.ReloadFile("desktop")
		v.ReloadFile("mobile")
		return
	}
	v.ReloadFile(v.source.layout)
	return
}

func (v *Visitor) ReloadFile(layout string) {
	cmd := ""
	switch v.fileType {
	case "css", "js", "img", "tmpl":
		cmd = v.fileType
	default:
		fmt.Println("WARNING: UNRECOGNIZED FILE TYPE", v.fileType)
	}
	
	source1 := v.source.JoinWithLayout("base")
	source2 := v.source.JoinWithLayout(layout)
	dest := v.source.JoinWithRootLayoutAndSubfolder(v.dest.root, layout, v.dest.subfolder)
	
	dlog.Println("Calculated destination as:", dest)
	
	proc, err := osutil.RunWithEnv(
		"framework/merge-handlers/" + cmd,
		nil,
		[]string{
			"WFDR_SOURCE_1=" + source1,
			"WFDR_SOURCE_2=" + source2,
			"WFDR_DEST=" + dest,
		})
	if err != nil {
		log.Println("Failed to run command", cmd, ":", err)
		return
	}
	proc.Wait()
}