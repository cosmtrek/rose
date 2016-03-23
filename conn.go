package main

import (
	"net"

	"github.com/cosmtrek/rose/protocol"
)

func connWrite(conn *net.Conn, msg string) {
	(*conn).Write(protocol.Pack([]byte(msg)))
}
