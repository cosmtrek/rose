package main

import (
	"github.com/cosmtrek/rose/protocol"
	"net"
)

func connWrite(conn *net.Conn, msg string) {
	(*conn).Write(protocol.Pack([]byte(msg)))
}
