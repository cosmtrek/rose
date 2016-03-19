package main

import (
	"log"
	"net"
	"time"

	"encoding/json"
	"github.com/cosmtrek/rose/protocol"
	"strconv"
)

const (
	ServerHost        = "localhost"
	ServerPort        = "3333"
	ConnType          = "tcp"
	LongSocketTimeout = 300
)

var (
	// TODO: need mutex
	OnlineUsers map[int]*net.Conn
	// messages
	ServerMessage = map[string]string{
		"existed":         time.Now().Format(time.RFC822) + " you're forced to exit",
		"unknown_actions": "unknown actions",
		"push_done":       "successfully push messages to all online users",
	}
)

func main() {
	ln, err := net.Listen(ConnType, ServerHost+":"+ServerPort)
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
	log.Println("Listening client", conn.RemoteAddr().String())
}

func heartbeating(conn net.Conn, message <-chan []byte, done chan<- bool) {
	for {
		select {
		case content := <-message:
			p := RequestParams{}
			if err := parseRequest(content, &p); err != nil {
				log.Println("request params broken")
				break
			}
			log.Println("Client " + conn.RemoteAddr().String() + " params - " + p.String())
			if p.Action == "ping" {
				if c, ok := OnlineUsers[p.Id]; ok {
					log.Printf("Found existed user: " + strconv.Itoa(p.Id))
					(*c).Write([]byte(ServerMessage["existed"]))
					(*c).Close()
					delete(OnlineUsers, p.Id)
				}

				OnlineUsers[p.Id] = &conn
				conn.SetDeadline(time.Now().Add(time.Duration(LongSocketTimeout) * time.Second))
			} else if p.Action == "push" {
				pushNotification(&OnlineUsers, []byte(p.Args))
				conn.Write([]byte(ServerMessage["push_done"]))
				conn.Close()
				done <- true
			} else {
				conn.Write([]byte(ServerMessage["unknown_actions"]))
				conn.Close()
				done <- true
			}
		case <-time.After(LongSocketTimeout * time.Second):
			log.Println("Client " + conn.RemoteAddr().String() + " exit")
			conn.Write([]byte("Closing..."))
			conn.Close()
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
			conn.Close()
			return
		default:
			n, _ := conn.Read(buf)
			tmpBuf = protocol.Unpack(append(tmpBuf, buf[:n]...), message)
		}
	}
}

type RequestParams struct {
	Id     int    `json:"id"`
	Action string `json:"action"`
	Args   string `json:"args,omitempty"`
}

func parseRequest(message []byte, p *RequestParams) error {
	if err := json.Unmarshal(message, p); err != nil {
		log.Println("json unmarshal error")
		return err
	}
	return nil
}

func (p *RequestParams) String() string {
	s := "id:" + strconv.Itoa(p.Id) + ", action:" + p.Action
	if p.Args != "" {
		s += ", args: " + p.Args
	}
	return s
}

func pushNotification(conns *map[int]*net.Conn, message []byte) {
	for _, c := range *conns {
		(*c).Write(message)
	}
}
