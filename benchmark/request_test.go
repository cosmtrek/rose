package main

import (
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/cosmtrek/rose/protocol"
)

// Just for fun :)
//   go test -bench .
//
// func BenchmarkHello(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		fmt.Println("hello")
//	}
// }

// TODO: need to improve
func BenchmarkConnectSequentially(b *testing.B) {
	server := "127.0.0.1:3333"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Println("Failed to resolve server")
		return
	}

	for i := 0; i < b.N; i++ {
		conn, err := setupConn(tcpAddr)
		if err != nil {
			fmt.Println(err)
		}
		done := make(chan bool)
		go sender(conn, i, done)
		<-done
	}
}

func BenchmarkConnectConcurrently(b *testing.B) {
	server := "127.0.0.1:3333"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Println("Failed to resolve server")
		return
	}

	b.RunParallel(func(pb *testing.PB) {
		i := 1
		for pb.Next() {
			conn, err := setupConn(tcpAddr)
			if err != nil {
				fmt.Println(err)
			}
			done := make(chan bool)
			i++
			go sender(conn, i, done)
			<-done
		}
	})
}

func setupConn(tcpAddr *net.TCPAddr) (net.Conn, error) {
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("Failed to connect server")
		return nil, err
	}
	return conn, nil
}

func sender(conn net.Conn, id int, done chan<- bool) {
	idS := strconv.Itoa(id)
	msg := "{\"id\":" + idS + ",\"action\":\"ping\"}"
	_, err := conn.Write(protocol.Pack([]byte(msg)))
	if err != nil {
		fmt.Println("Client " + idS + " cannot write to remote connection")
		return
	}
	done <- true
	// Not close conn
}
