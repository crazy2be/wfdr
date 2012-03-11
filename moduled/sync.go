package moduled

import (
	"errors"
	"exp/inotify"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"wfdr/pathbits"

	"github.com/crazy2be/osutil"
)

// Sync recursively traverses the directory tree for source, updating the files in dest if the source files have been modified more recently than the destination files.
func (cm *CacheMonitor) Sync() error {
	err := cm.syncDir(cm.source.Name(), cm.dest.Name())
	if err != nil {
		return err
	}
	return nil
}

// syncDir is called once, recursively, for every directory in source.
func (cm *CacheMonitor) syncDir(source string, dest string) error {
	sourcef, err := os.Open(source)
	if err != nil {
		return err
	}
	fnames, err := sourcef.Readdirnames(-1)
	if err != nil {
		return err
	}
	sort.Strings(fnames)

	basefile := ""
	genlays := make([]bool, len(knownLayouts))
	for i := range fnames {
		fname := fnames[i]
		fi, err := os.Stat(path.Join(source, fname))
		if err != nil {
			return err
		}
		if fi.IsDir() {
			err = os.Mkdir(path.Join(dest, fname), 0744)
			cm.syncDir(path.Join(source, fname), path.Join(dest, fname))
			continue
		}

		// Check if this file is layout specific. If so, process it, then move onto the next file.
		layout := cm.layoutSpecific(fname)
		if layout != -1 {
			if !strings.HasPrefix(pathbits.Name(fname), pathbits.Name(basefile)) {
				return errors.New(fmt.Sprintf("wfdr/moduled: File %s has no complementry default file. (i.e. foobar_<layout> exists with no foobar file to complement it)", fname))
			}
			err = cm.updateFile(
				path.Join(source, basefile),
				path.Join(source, fname),
				path.Join(dest, knownLayouts[layout], basefile))
			if err != nil {
				return err
			}
			genlays[layout] = true
			continue
		}

		// Generate files for any layouts that do not have custom files.
		err = cm.genRemaining(genlays, source, dest, basefile)
		if err != nil {
			return err
		}

		basefile = fname
	}
	err = cm.genRemaining(genlays, source, dest, basefile)
	if err != nil {
		return err
	}
	return nil
}

func (cm *CacheMonitor) Deamon(errors chan<- error) {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		errors <- err
		return
	}
	err = cm.daemonInit(watcher, cm.source.Name())
	for {
		select {
		case ev := <-watcher.Event:
			cm.daemonEvent(ev)
		case err := <-watcher.Error:
			errors <- err
		}
	}
}

// daemonInit walks the directory tree rooted at dir, recursively adding each subdirectory of dir to watcher's watchlist.
func (cm *CacheMonitor) daemonInit(watcher *inotify.Watcher, dir string) error {
	err := watcher.AddWatch(dir, 0)
	if err != nil {
		return err
	}
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	fis, err := f.Readdir(-1)
	for i := range fis {
		fi := fis[i]
		if fi.IsDir() {
			err = cm.daemonInit(watcher, path.Join(dir, fi.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (cm *CacheMonitor) daemonEvent(ev *inotify.Event) {
	// Ignore qt temp files (from kate)
	if strings.Index(strings.ToLower(ev.Name), "qt_temp") != -1 {
		cm.dlog.Println("Ignoring kate temp file...")
		return
	}
	// Ignore temp/backup files.
	if strings.HasSuffix(ev.Name, "~") {
		cm.dlog.Println("Ignoring backup file (suffix ~)")
		return
	}
}

func (cm *CacheMonitor) genRemaining(genlays []bool, source, dest, name string) error {
	// Can't generate for empty file name...
	if name == "" {
		cm.dlog.Println("genRemaining called with empty name (typically happens for the first file in each directory)")
		return nil
	}
	for i := range genlays {
		if genlays[i] == true {
			genlays[i] = false
			continue
		}
		err := cm.updateFile(path.Join(source, name), "", path.Join(dest, knownLayouts[i], name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (cm *CacheMonitor) layoutSpecific(name string) int {
	for i := range knownLayouts {
		layout := knownLayouts[i]

		matched, err := path.Match("*_"+layout+".*", name)
		if err != nil {
			// Only ever happens when our pattern syntax is incorrect, in which case panicing seems like the best solution.
			panic(err)
		}

		if matched {
			return i
		}
	}
	return -1
}

func (cm *CacheMonitor) updateFile(source1, source2, dest string) error {
	fi1, err := os.Stat(source1)
	if err != nil {
		return err
	}
	destfi, err := os.Stat(dest)
	// We want to reload if either
	// 1) The source was modified after the destination, or
	// 2) The destination file does not exist (and thus should be created).
	if err != nil || fi1.ModTime().After(destfi.ModTime()) {
		return cm.reloadFile(source1, source2, dest)
	}

	if source2 != "" {
		fi2, err := os.Stat(source2)
		if err != nil {
			return err
		}
		if fi2.ModTime().After(destfi.ModTime()) {
			return cm.reloadFile(source1, source2, dest)
		}
	}
	return nil
}

func (cm *CacheMonitor) reloadFile(source1, source2, dest string) error {
	// Create the destination directory if it does not exist.
	fpath, _ := path.Split(dest)
	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(fpath, 0744)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	proc, err := osutil.RunWithEnv(
		"framework/merge-handlers/"+cm.typ,
		nil,
		[]string{
			"WFDR_SOURCE_1=" + source1,
			"WFDR_SOURCE_2=" + source2,
			"WFDR_DEST=" + dest,
		})
	if err != nil {
		return errors.New(fmt.Sprintln("Failed to run command", cm.typ, ":", err))
	}
	proc.Wait()
	return nil
}
