package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"syscall"
	"time"
	// Local imports
	"github.com/crazy2be/osutil"
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

func exitWait(process *os.Process, done chan<- error) {
	for {
		ps, err := process.Wait()
		if err != nil {
			done <- err
			return
		}
		if ps.Exited() {
			done <- nil
			return
		}
	}
}

// Waits for a process with the given PID to exit.
func ExitWait(process *os.Process) error {
	done := make(chan error)
	go exitWait(process, done)
	select {
	case err := <-done:
		return err
	case <-time.After(15 * time.Second):
		return errors.New(fmt.Sprintf("Timed out waiting for process with PID %d to exit. You might want to clean it up manually.", process.Pid))
	}
	panic("Not reached!")
}

func (m *Module) Stop() error {
	log.Printf("Stopping module %s", m.Name)
	if m.Name == "internal:shared" {
		err := m.SyncProcess.Signal(os.Interrupt)
		if err != nil {
			return err
		}
		err = ExitWait(m.SyncProcess)
		if err != nil {
			return err
		}
		delete(modules, m.Name)
		return nil
	}
	// 0x02, or SIGINT.
	err1 := m.MainProcess.Signal(os.Interrupt)
	err2 := m.SyncProcess.Signal(os.Interrupt)
	if err1 != nil {
		return errors.New(fmt.Sprintf("Failed to stop module %s:", m.Name, err1))
	}
	if err2 != nil {
		return errors.New(fmt.Sprintf("Failed to stop sync process for module %s (PID %d), you should probably stop it manually.", m.Name, m.SyncProcess.Pid))
	}

	err1 = ExitWait(m.MainProcess)
	if err1 != nil {
		return err1
	}

	err2 = ExitWait(m.SyncProcess)
	if err2 != nil {
		return err2
	}

	delete(modules, m.Name)
	return nil
}

func (m *Module) IsRunning() bool {
	pid := 0
	if m.Name == "internal:shared" {
		pid = m.SyncProcess.Pid
	} else {
		pid = m.MainProcess.Pid
	}

	var waitstatus syscall.WaitStatus
	wpid, err := syscall.Wait4(pid, &waitstatus, syscall.WNOHANG|syscall.WUNTRACED, nil)
	if err != nil {
		// When would this happen?
		log.Println("Unable to get process wait status:", err)
		// Assume it is not running
		return false
	}

	// If status is not available, the pid is 0.
	if wpid == 0 {
		return true
	}
	if waitstatus.Exited() {
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

	log.Printf("Syncing shared resources...")
	var err error
	_, err = osutil.WaitRun("shared-init", nil)
	if err != nil {
		return nil, err
	}
	log.Println("Done syncing shared resources.")

	mod := new(Module)
	mod.Name = name

	cmd, err := osutil.Run("shared-daemon", nil)
	if err != nil {
		return nil, err
	}

	mod.SyncProcess = cmd.Process
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

	err := JailInit(moddir, jaildir, name)
	if err != nil {
		return nil, err
	}

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
