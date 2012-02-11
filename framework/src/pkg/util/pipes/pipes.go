// A few useful functions for dealing with named pipes, and using them as two-way connections.
package pipes

import (
	"errors"
	"fmt"
	"os"
	"time"
	// Local imports
	"github.com/crazy2be/osutil"
)

type PipeReadWriteCloser struct {
	infile  *os.File
	outfile *os.File
}

func NewPipeReadWriteCloser(infile, outfile *os.File) *PipeReadWriteCloser {
	return &PipeReadWriteCloser{infile, outfile}
}

// Creates the pipe if it does not exist, and then opens it using OpenPipe().
func MakeAndOpen(pipename string) (*os.File, error) {
	if !osutil.FileExists(pipename) {
		Mkfifo(pipename, 0644)
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

	file, err := os.OpenFile(pipename, os.O_RDWR|os.O_NONBLOCK, 0644)
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error opening pipe ", pipename, ": ", err))
	}
	return file, nil
}

// Implements a semi-blocking interface to the non-blocking underlying pipe (assuming the pipe was opened with CheckPipe())
// BUG: Slower than a normal read, because it Sleep()s in a loop polling for data. 
func (prwc *PipeReadWriteCloser) Read(buf []byte) (int, error) {
	n, err := prwc.infile.Read(buf)
	for n == 0 {
		n, err = prwc.infile.Read(buf)
		time.Sleep(1000 * 1000 * 10) //  1/100th of a second, prevents it from using 100% of the CPU.
	}
	return n, err
}

func (prwc *PipeReadWriteCloser) Write(buf []byte) (int, error) {
	return prwc.outfile.Write(buf)
}

func (prwc *PipeReadWriteCloser) Close() error {
	err := prwc.infile.Close()
	if err != nil {
		return err
	}
	err = prwc.outfile.Close()
	if err != nil {
		return err
	}
	return nil
}

// TODO: Use syscall.Mkfifo, but it doesn't seem to work atm.
// NOTE: perms are ignored, since this calls the mkfifo program to do it's work.
// WARNING: Hackish implementation.
func Mkfifo(name string, perms uint32) {
	osutil.WaitRun("mkfifo", []string{name})
}
