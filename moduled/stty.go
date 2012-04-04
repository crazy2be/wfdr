package moduled

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

type sttyState string

func (stty sttyState) Undo() {
	proc := exec.Command("stty", string(stty))
	proc.Stdin = os.Stdin
	err := proc.Run()
	if err != nil {
		log.Println("Warning: Failed to reset terminal options with to " + stty)
		return
	}
	return
}

func readStty() (sttyState, error) {
	proc := exec.Command("/bin/stty", "-g")
	proc.Stdin = os.Stdin
	output, err := proc.Output()
	if err != nil {
		return sttyState(output), errors.New("Reading tty state with 'stty -g': " + err.Error())
	}
	return sttyState(output), nil
}

func SttyCbreak() (sttyState, error) {
	state, err := readStty()
	if err != nil {
		return state, err
	}
	proc := exec.Command("/bin/stty", "cbreak")
	proc.Stdin = os.Stdin
	err = proc.Run()
	if err != nil {
		return state, errors.New("Executing 'stty raw': " + err.Error())
	}
	return state, nil
}
