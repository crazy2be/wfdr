// Provides functions for communicating with the module deamon through the rpc module.
package moduled

import (
	"net/rpc"
	"net/rpc/jsonrpc"
	// Local imports
	"wfdr/pipes"
)

// Connects to the pipe files, in order to allow this program to sent commands to the process management deamon.
func RPCConnect() (*rpc.Client, error) {
	// Pipes are reversed from what you would expect because we are connecting as a client, and they are named based on how the server uses them. Thus, the out pipe for the server is the in pipe for us.
	outpipe := "cache/wfdr-deamon-pipe-in"
	inpipe := "cache/wfdr-deamon-pipe-out"

	infile, err := pipes.Open(inpipe)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	outfile, err := pipes.Open(outpipe)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}

	rwc := &pipes.PipeReadWriteCloser{infile, outfile}

	return jsonrpc.NewClient(rwc), nil
}
