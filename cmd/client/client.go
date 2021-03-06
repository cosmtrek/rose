package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/cosmtrek/rose/protocol"
)

func sender(conn net.Conn, id int, done chan<- bool) {
	defer func() {
		done <- true
	}()

	msg := "{\"id\":" + strconv.Itoa(id) + ",\"action\":\"ping\"}"
	if _, err := conn.Write(protocol.Pack([]byte(msg))); err != nil {
		fmt.Println("Cannot write to remote connection and exiting...")
		os.Exit(1)
	}
}

func reader(conn net.Conn, done chan bool) {
	defer func() {
		done <- true
	}()

	buf := make([]byte, 1024)
	tmpBuf := make([]byte, 1024)
	message := make(chan []byte)

	go func(message chan []byte, done chan bool) {
		for {
			select {
			case m := <-message:
				log.Println(string(m))
			case <-done:
				return
			}
		}
	}(message, done)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Cannot read from remote connection and exiting...")
			os.Exit(1)
		}

		tmpBuf = protocol.Unpack(append(tmpBuf, buf[:n]...), message)
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

	id := flag.Int("id", 1, "client id")
	flag.Parse()

	var sendDone chan bool
	var readDone chan bool

	go sender(conn, *id, sendDone)
	go reader(conn, readDone)

	<-sendDone
	<-readDone
}
