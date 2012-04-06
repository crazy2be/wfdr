package moduled

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

type sttyState string

// Undo restores the terminal to whichever state it was in prior to a call to SttyCbreak(). See the documentation for SttyCbreak for more details on proper usage.
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

// SttyCbreak puts the connected terminal, if any, into cbreak mode (cooked mode but with escape sequences still having their usual effect). Returns the previous state of the terminal, and an error, if any. The effects of calling SttyCbreak() should be undone to return the terminal to a sane state upon program exit, as shown below:
//	state, err := SttyCbreak()
//	defer state.Undo()
//	if err != nil {
//		// Do something with error
//	}
//	// Terminal is now in cbreak mode!
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
