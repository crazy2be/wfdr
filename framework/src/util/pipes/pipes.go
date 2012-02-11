// A few useful functions for dealing with named pipes, and using them as two-way connections.
package pipes

import (
	"syscall"
	"errors"
	"fmt"
	"os"
	// Local imports
	"github.com/crazy2be/osutil"
)

type PipeReadWriteCloser struct {
	Input  *os.File
	Output *os.File
}

// Creates the pipe if it does not exist, and then opens it using OpenPipe().
func MakeAndOpen(pipename string) (*os.File, error) {
	if !osutil.FileExists(pipename) {
		syscall.Mkfifo(pipename, 0644)
	}
	return Open(pipename)
}

// Checks a named pipe for sanity, ensuring that the caller can read the pipe, and that it is, in fact, a pipe. Returns the file if it could be successfully opened.
func Open(pipename string) (*os.File, error) {
	finfo, err := os.Stat(pipename)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error stating pipe ", pipename, ": ", err))
	}

	if finfo.Mode()&os.ModeNamedPipe != os.ModeNamedPipe {
		return nil, errors.New(fmt.Sprint("Specified pipe file ", pipename, " is not a named pipe type file! Remove it to have a pipe file created in it's place."))
	}

	file, err := os.OpenFile(pipename, os.O_RDWR, 0644)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error opening pipe ", pipename, ": ", err))
	}
	return file, nil
}

// Read reads from the input pipe, and is effectively equivelent to reading from prwc.Input directly.
func (prwc *PipeReadWriteCloser) Read(buf []byte) (int, error) {
	return prwc.Input.Read(buf)
}

func (prwc *PipeReadWriteCloser) Write(buf []byte) (int, error) {
	return prwc.Output.Write(buf)
}

// Close closes both the input and the output pipes used by the prwc, returning an error if either operation fails.
func (prwc *PipeReadWriteCloser) Close() error {
	err := prwc.Input.Close()
	if err != nil {
		return err
	}
	err = prwc.Output.Close()
	if err != nil {
		return err
	}
	return nil
}