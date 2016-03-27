package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/cosmtrek/rose/protocol"
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
	go guardSignal()

	ln, err := net.Listen("tcp", config.ServerHost+":"+config.ServerPort)
	if err != nil {
		Fatal(err)
	}
	defer ln.Close()

	go processRequest()

	for {
		conn, err := ln.Accept()
		if err != nil {
			errl.Println(err)
			break
		}
		go handleRequest(conn)
	}
}

func processRequest() {
	for {
		select {
		case rp := <-RequestQueue:
			c := *rp.Conn
			r := rp.Message
			info.Println(c.RemoteAddr().String() + " " + r.String())
			if r.Action == Ping {
				if c, ok := global.getOnlineUser(r.Id); ok {
					debug.Println("Found existed user " + strconv.Itoa(r.Id))
					p := newResponse("existed", ResponsePush)
					connWrite(c, p.Json())
					(*c).Close()
					global.deleteOnlineUser(r.Id)
				}

				global.addOnlineUser(r.Id, rp.Conn)
				c.SetDeadline(time.Now().Add(time.Duration(config.SocketTimeout) * time.Second))
				p := newActionResponse(Ping)
				connWrite(&c, p.Json())
			} else if r.Action == Push {
				pushMessage([]byte(r.Args))
				p := newActionResponse(Push)
				connWrite(&c, p.Json())
			} else {
				p := newResponse("unknown_actions", ResponseReply)
				connWrite(&c, p.Json())
			}
		}
	}
}

type RequestPackage struct {
	Conn    *net.Conn
	Message Request
}

var (
	RequestQueue = make(chan RequestPackage, 1000)
)

func systemStat() {
	debug.Printf("Current request queue size: %d", len(RequestQueue))
	debug.Printf("OnlineUsers map size: %d", len(global.OnlineUsers))
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 64)
	tmpBuf := make([]byte, 64)
	m := make(chan []byte, 1)

	for {
		select {
		case r := <-m:
			p := RequestPackage{
				Conn:    &conn,
				Message: Request{},
			}
			if err := parseRequest(r, &p.Message); err != nil {
				errl.Println("Failed to parse request params")
				break
			}
			debug.Printf("RequestPackage: %v", p)
			RequestQueue <- p
		default:
			n, err := conn.Read(buf)
			if err != nil {
				conn.Close()
				close(m)
				return
			}

			tmpBuf = protocol.Unpack(append(tmpBuf, buf[:n]...), m)
		}
	}
}

func pushMessage(message []byte) {
	for _, c := range global.OnlineUsers {
		p := newPushResponse(string(message))
		connWrite(c, p.Json())
	}
}
