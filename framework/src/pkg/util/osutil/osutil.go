// Utility functions for functionality you would expect to find in the os package, such as running programs and testing if files exist.
package osutil

import (
	"os"
	"log"
	"exec"
	"path"
	"strings"
)

// Checks if a file exists, returns true if the file exists and can be opened, false otherwise. If the file exists but cannot be opened, it still returns false.
// NOTE: Assumes that any error opening the file for reading means that the file does not exist. This is not always true, but is probably the best behaviour anyway. There is rarely a usage case where you want to check if a file exists, but not if you can open it.
func FileExists(name string) bool {
	_, err := os.Stat(name)
	
	if err != nil {
		return false
	}
	return true
}

func findEnv(env []string, name string) string {
	for i := 0; i < len(env); i++ {
		foundname, value := parseEnv(env[i])
		if foundname == name {
			return value
		}
	}
	return ""
}

func parseEnv(variable string) (string, string) {
	namevalue := strings.Split(variable, "=", 2)
	name := namevalue[0]
	value := ""
	if len(namevalue) > 1 {
		value = namevalue[1]
	}
	return name, value
}

func mergeEnv(base []string, overrides []string) []string {
	envmap := make(map[string]string, len(base))
	for i := 0; i < len(base); i++ {
		name, value := parseEnv(base[i])
		envmap[name] = value
	}
	for i := 0; i < len(overrides); i++ {
		name, value := parseEnv(overrides[i])
		envmap[name] = value
	}
	envslice := make([]string, len(envmap))
	i := 0
	for name, value := range envmap {
		envslice[i] = name + "=" + value
		i++
	}
	return envslice
}

func findCmd(PATH, cmd string) (string, os.Error) {
	// If the name contains a /, e.g. framework/foo/blar or ./foo, then try the path directly.
	if strings.Index(cmd, "/") != -1 {
		if FileExists(cmd) {
			return cmd, nil
		}
	}
	paths := strings.Split(PATH, ":", -1)
	for i := 0; i < len(paths); i++ {
		lookpath := paths[i]
		binpath := path.Join(lookpath, cmd)
		if FileExists(binpath) {
			return binpath, nil
		}
	}
	return "", os.NewError("Command " + cmd + " not found! Path is " + PATH)
}

func RunWithEnvAndWd(command string, args []string, env []string, wd string) (proc *exec.Cmd, err os.Error) {
	//log.Println(command, args)
	hho := exec.PassThrough
	args = prepend(args, command)
	env = mergeEnv(os.Environ(), env)
	
	binpath, err := findCmd(findEnv(env, "PATH"), command)
	if err != nil {
		return nil, err
	}
	
	proc, err = exec.Run(binpath, args, env, wd, hho, hho, hho)
	if err != nil {
		log.Print("Error running command ", command, ": ", err, "\n")
		return nil, err
	}
	return
}

// More advanced, runs a program with a custom enviroment. Note that the normal enviroment is also passed here, as that is what is typically desired. Enviroment is a slice of strings, with each string usually having the form "NAME=VALUE". If you pass an enviroment string of the form NAME=VALUE that has the same name as an existing enviroment string, your value overwrites the value of the other variable.
func RunWithEnv(command string, args []string, env []string) (proc *exec.Cmd, err os.Error) {
	return RunWithEnvAndWd(command, args, env, ".")
}

// Simple way to run most programs. Searches for the program in PATH, and runs the first found program. Args need not contain the program name as the zeroth argument, it is prepended automatically.
func Run(command string, args []string) (proc *exec.Cmd, err os.Error) {
	return RunWithEnv(command, args, []string{})
}

// Runs a command using Run(), but waits for it to complete before returning.
func WaitRun(command string, args []string) (proc *exec.Cmd, err os.Error) {
	proc, err = Run(command, args)
	if err != nil {
		return
	}
	proc.Close()
	return
}

// Doesn't really belong here. Oh well...
func prepend(orig []string, prep ...string) []string {
	//log.Println(orig, prep)
	arr := make([]string, len(prep)+len(orig))
	for i := len(prep); i < len(arr); i++ {
		arr[i] = orig[i-len(prep)]
	}
	for i := 0; i < len(prep); i++ {
		arr[i] = prep[i]
	}
	//log.Println(arr)
	return arr
}