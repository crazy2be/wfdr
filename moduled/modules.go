package moduled

import (
	"errors"
	"net/rpc"
)

// StartModule tells the rpc server given by conn to stop the module given by name.
func StartModule(conn *rpc.Client, name string) error {
	var dummy int = 1000
	err := conn.Call("ModuleSrv.Start", &name, &dummy)
	if err != nil {
		return errors.New("Error starting " + name + ":" + err.Error())
	}
	return nil
}

func StopModule(conn *rpc.Client, name string) error {
	var dummy int = 1000
	err := conn.Call("ModuleSrv.Stop", &name, &dummy)
	if err != nil {
		return errors.New("Error stopping " + name + ":" + err.Error())
	}
	return nil
}