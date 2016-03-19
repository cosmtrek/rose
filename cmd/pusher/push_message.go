package main

import (
	"fmt"
	"github.com/cosmtrek/rose/protocol"
	"log"
	"net"
	"os"
)

func sender(conn net.Conn, done chan<- bool) {
	defer func() {
		done <- true
	}()

	msg := "{\"id\":0,\"action\":\"push\", \"args\":\"what's up?\"}"
	if _, err := conn.Write(protocol.Pack([]byte(msg))); err != nil {
		fmt.Println("Cannot write to remote connection and exiting...")
		os.Exit(1)
	}
}

func reader(conn net.Conn, done chan<- bool) {
	defer func() {
		done <- true
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Cannot read from remote connection and exiting...")
			os.Exit(1)
		}
		log.Println(string(buf[:n]))
	}
}

func main() {
	server := "127.0.0.1:3333"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Connect successfully")

	var sendDone chan bool
	var readDone chan bool

	go sender(conn, sendDone)
	go reader(conn, readDone)

	<-sendDone
	<-readDone
}
