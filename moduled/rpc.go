package moduled

import (
	"errors"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Conn struct {
	client *rpc.Client
}

// Connects to the pipe files, in order to allow this program to sent commands to the process management deamon.
func RPCConnect() (*Conn, error) {
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

	return &Conn{client: jsonrpc.NewClient(rwc)}, nil
}

func (c *Conn) Start(name string) error {
	var dummy int = 1000
	err := c.client.Call("ModuleSrv.Start", &name, &dummy)
	if err != nil {
		return errors.New("Error starting " + name + ": " + err.Error())
	}
	return nil
}

func (c *Conn) Stop(name string) error {
	var dummy int = 1000
	err := c.client.Call("ModuleSrv.Stop", &name, &dummy)
	if err != nil {
		return errors.New("Error stopping " + name + ": " + err.Error())
	}
	return nil
}

func (c *Conn) Restart(name string) error {
	err := c.Stop(name)
	if err != nil {
		return err
	}
	return c.Start(name)
}

func (c *Conn) Status(name string) (running bool, err error) {
	err = c.client.Call("ModuleSrv.Status", &name, &running)
	if err != nil {
		return false, errors.New("Error getting status for " + name + ": " + err.Error())
	}
	return running, nil
}

func (c *Conn) Close() error {
	return c.client.Close()
}

func (c *Conn) InterpretCommand(args []string) error {
	if args == nil || len(args) < 1 {
		return errors.New("No command provided!")
	}
	cmd := args[0]
	args = args[1:]
	
	switch (cmd) {
	case "start":
		for _, module := range args {
			err := c.Start(module)
			if err != nil {
				return err
			}
		}
	case "stop":
		for _, module := range args {
			if len(args) < 1 {
				return errors.New("Module name required!")
			}
			return c.Stop(module)
		}
	case "restart":
		for _, module := range args {
			if len(args) < 1 {
				return errors.New("Module name required!")
			}
			return c.Restart(module)
		}
	}
	
	return nil
}