package moduled

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	// Local imports
	"github.com/crazy2be/osutil"
)

type PipeReadWriteCloser struct {
	Input  *os.File
	Output *os.File
}

// Creates the pipe if it does not exist, and then opens it. If the file exists, it ensures that the existing file is a pipe. Returns the file if it could be successfully opened.
func OpenPipe(pipename string) (*os.File, error) {
	if !osutil.FileExists(pipename) {
		syscall.Mkfifo(pipename, 0644)
	}

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
	err1 := prwc.Input.Close()
	err2 := prwc.Output.Close()
	// Should attempt to close both pipes before returning any error.
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}
