package moduled

import (
	"net/rpc"
	"io"
)

type Shell struct {
	conn *rpc.Client
	rd   io.Reader
}

func NewShell(conn *rpc.Client, rd io.Reader) {
	s := new(Shell)
	s.conn = conn
	s.rd = rd
}