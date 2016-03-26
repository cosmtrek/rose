package main

import (
	"net"

	"github.com/cosmtrek/rose/protocol"
)

func connWrite(conn *net.Conn, msg string) {
	if _, err := (*conn).Write(protocol.Pack([]byte(msg))); err != nil {
		errl.Println(err)
		(*conn).Close()
	}
}
