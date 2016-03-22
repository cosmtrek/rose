package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/cosmtrek/rose/protocol"
)

var (
	// TODO: need mutex
	OnlineUsers map[int]*net.Conn
)

func main() {
	fmt.Println(`
  _____
 |  __ \
 | |__) |   ___    ___    ___
 |  _  /   / _ \  / __|  / _ \
 | | \ \  | (_) | \__ \ |  __/
 |_|  \_\  \___/  |___/  \___|
`)

	initConfig()

	ln, err := net.Listen("tcp", config.ServerHost+":"+config.ServerPort)
	CheckErr(err)
	defer ln.Close()

	OnlineUsers = make(map[int]*net.Conn)

	for {
		conn, err := ln.Accept()
		CheckErr(err)
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	message := make(chan []byte, 64)
	done := make(chan bool)

	go readRequest(conn, message, done)
	go heartbeating(conn, message, done)
}

func heartbeating(conn net.Conn, message <-chan []byte, done chan<- bool) {
	for {
		select {
		case content := <-message:
			r := Request{}
			if err := parseRequest(content, &r); err != nil {
				errl.Println("Failed to parse request params")
				break
			}
			debug.Println("Client " + conn.RemoteAddr().String() + " " + r.String())
			if r.Action == Ping {
				if c, ok := OnlineUsers[r.Id]; ok {
					debug.Println("Found existed user " + strconv.Itoa(r.Id))
					p := newResponse("existed", ResponsePush)
					connWrite(c, p.Json())
					(*c).Close()
					delete(OnlineUsers, r.Id)
				}

				OnlineUsers[r.Id] = &conn
				conn.SetDeadline(time.Now().Add(time.Duration(config.SocketTimeout) * time.Second))
				p := newActionResponse(Ping)
				connWrite(&conn, p.Json())
			} else if r.Action == Push {
				pushMessage(&OnlineUsers, []byte(r.Args))
				p := newActionResponse(Push)
				connWrite(&conn, p.Json())
				done <- true
				return
			} else {
				p := newResponse("unknown_actions", ResponseReply)
				connWrite(&conn, p.Json())
				done <- true
				return
			}
		case <-time.After(time.Duration(config.SocketTimeout) * time.Second):
			debug.Println("Client " + conn.RemoteAddr().String() + " exit")
			done <- true
			return
		}
	}
}

func readRequest(conn net.Conn, message chan<- []byte, done <-chan bool) {
	buf := make([]byte, 1024)
	tmpBuf := make([]byte, 1024)

	for {
		select {
		case <-done:
			go updateOnlineUsers(&conn)
			conn.Close()
			return
		default:
			n, _ := conn.Read(buf)
			tmpBuf = protocol.Unpack(append(tmpBuf, buf[:n]...), message)
		}
	}
}

func pushMessage(conns *map[int]*net.Conn, message []byte) {
	for _, c := range *conns {
		p := newPushResponse(string(message))
		connWrite(c, p.Json())
	}
}

func updateOnlineUsers(conn *net.Conn) {
	debug.Println("Updating online users...")
	for k, v := range OnlineUsers {
		if *v == *conn {
			debug.Println("Delete user " + strconv.Itoa(k))
			delete(OnlineUsers, k)
		}
	}
}
