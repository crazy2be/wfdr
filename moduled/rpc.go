package moduled

import (
	"errors"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// Connects to the pipe files, in order to allow this program to sent commands to the process management deamon.
func RPCConnect() (*rpc.Client, error) {
	// Pipes are reversed from what you would expect because we are connecting as a client, and they are named based on how the server uses them. Thus, the out pipe for the server is the in pipe for us.
	outpipe := "cache/wfdr-deamon-pipe-in"
	inpipe := "cache/wfdr-deamon-pipe-out"

	infile, err := OpenPipe(inpipe)
	if err != nil {
		return nil, errors.New(err.Error() + ". In all likelyhood, wfdr-deamon is not running. Start it and this error should go away! It's also possible that you don't have permission to talk to the deamon process, in which case i can't help you.")
	}
	outfile, err := OpenPipe(outpipe)
	if err != nil {
		return nil, err
	}

	rwc := &PipeReadWriteCloser{Input: infile, Output: outfile}

	return jsonrpc.NewClient(rwc), nil
}
