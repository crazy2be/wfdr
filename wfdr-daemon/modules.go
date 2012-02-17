package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"syscall"
	// Local imports
	"github.com/crazy2be/osutil"
	"util/dlog"
)

type Module struct {
	Name string
	// The process of the actual module
	MainProcess *os.Process
	// The process of the sync program, updates css and js automatically when changed (on supported systems)
	SyncProcess *os.Process
}

var modules = make(map[string]*Module)

func GetModule(name string) (*Module, error) {
	module, exists := modules[name]
	if !exists {
		return nil, errors.New("Module is not started!")
	}
	return module, nil
}

func ModuleRunning(name string) bool {
	module, _ := modules[name]
	if module == nil {
		return false
	}
	return module.IsRunning()
}

func (m *Module) Stop() error {
	// 0x02, or SIGINT.
	syscall.Kill(m.MainProcess.Pid, 0x02)
	syscall.Kill(m.SyncProcess.Pid, 0x02)
	// TODO: Check for syscall errors
	delete(modules, m.Name)
	return nil
}

func (m *Module) IsRunning() bool {
	//dlog.Printf("%#v %#v %#v\n", m, m.MainProcess, m.SyncProcess)
	waitmsg, err := os.Wait(m.MainProcess.Pid, os.WNOHANG|os.WUNTRACED)
	if err != nil {
		// TODO: When would this happen?
		dlog.Println("Unable to get process wait status:", err)
	}
	//dlog.Printf("%#v\n", waitmsg)
	// If status is not available, the pid is 0.
	if waitmsg.Pid == 0 {
		return true
	}
	if waitmsg.WaitStatus.Exited() {
		delete(modules, m.Name)
		return false
	}
	return true
}

func JailInit(moddir, jaildir, modname string) error {
	osutil.WaitRun("wfdr-reload-shared", nil)
	setup, err := osutil.RunWithEnv("jail-init", nil, []string{"WFDR_MODDIR=" + moddir, "WFDR_JAILDIR=" + jaildir, "WFDR_MODNAME=" + modname})
	if err != nil {
		return errors.New(fmt.Sprint("Could not run script to initialize jail:", err, " PATH:", os.Getenv("PATH")))
	}
	setup.Wait()
	return nil
}

func StartSync(moddir, jaildir, modname string) (*os.Process, error) {
	//hho := exec.PassThrough
	deamon, err := osutil.RunWithEnv("jail-deamon", nil, []string{"WFDR_MODDIR=" + moddir, "WFDR_JAILDIR=" + jaildir, "WFDR_MODNAME=" + modname})
	if err != nil {
		return nil, errors.New(fmt.Sprint("Could not start sync deamon, css, js, and template files will not be synced:", err))
	}
	return deamon.Process, nil
}

// Syncronizes the shared resources and starts the deamon to sync them.
func StartSharedSync() (*Module, error) {
	name := "internal:shared"
	if ModuleRunning(name) {
		return GetModule(name)
	}
	var err error
	_, err = osutil.WaitRun("shared-init", nil)
	if err != nil {
		return nil, err
	}
	mod := new(Module)
	mod.Name = name
	cmd, err := osutil.Run("shared-daemon", nil)
	if err != nil {
		return nil, err
	}
	mod.SyncProcess = cmd.Process
	mod.MainProcess = cmd.Process
	modules[name] = mod
	return mod, nil
}

func StartModule(name string) (*Module, error) {
	if ModuleRunning(name) {
		return nil, errors.New(fmt.Sprint("The module seems to be already started..."))
	}

	// Handle special/internal modules
	switch name {
	case "internal:shared":
		return StartSharedSync()
	}

	cwd, _ := os.Getwd()
	jaildir := path.Join(cwd, "jails/"+name)
	moddir := path.Join(cwd, "modules/"+name)
	JailInit(moddir, jaildir, name)

	path := jaildir + "/sh:" + jaildir + "/bin:" + os.Getenv("PATH")

	if !osutil.FileExists(jaildir + "/sh/run") {
		cp, err := osutil.WaitRun("cp", []string{"framework/sh/jail-run", jaildir + "/sh/run"})
		if err != nil {
			return nil, errors.New(fmt.Sprint("Error copying default run file, cannot continue:", err))
		}
		cp.Wait()
	}

	modulep, err := osutil.RunWithEnvAndWd("run", []string{name}, []string{"PATH=" + path}, jaildir)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Could not start module "+name+"!:", err))
	}
	pid := modulep.Process.Pid

	module := new(Module)
	module.Name = name
	module.MainProcess = modulep.Process
	// wtf? [sic]
	defer func() {
		modules[name] = module
	}()

	log.Println("Started Module "+name+"! PID:", pid)

	// Start the sync deamon, syncronizes css, js, and templates in the background
	syncproc, err := StartSync(moddir, jaildir, name)
	if err != nil {
		return module, err
	}

	module.SyncProcess = syncproc

	return module, nil
}

func StopModule(name string) error {
	if !ModuleRunning(name) {
		return errors.New("The module cannot possibly be stopped, as it does not appear to be running.")
	}

	module, err := GetModule(name)
	if err != nil {
		return err
	}
	return module.Stop()
}
