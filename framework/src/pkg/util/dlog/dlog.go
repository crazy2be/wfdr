// A simple debug logging package that provides facilities for logging things to a logfile.
package dlog

import (
	"os"
	"fmt"
	"log"
)

var dlog = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

// Logs to stdout, which will end up in a log file.
func Println(v ...interface{}) {
	dlog.Output(2, fmt.Sprintln(v...))
}

// Logs to stdout, which will end up in a log file.
func Printf(format string, v ...interface{}) {
	dlog.Output(2, fmt.Sprintf(format, v...))
}

// Outputs to stderr AND stdout
func Fatalln(v ...interface{}) {
	dlog.Output(2, fmt.Sprintln(v...))
	fmt.Fprintln(os.Stderr, v...)
}

// Output to stderr AND stdout
func Fatalf(format string, v ...interface{}) {
	dlog.Output(2, fmt.Sprintf(format, v...))
	fmt.Fprintf(os.Stderr, format, v...)
}