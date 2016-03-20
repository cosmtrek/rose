package main

import (
	"encoding/json"
	"log"
	"strconv"
)

const (
	Ping = "ping"
	Push = "push"
)

var (
	ReqResp = map[string]string{
		"ping": "pong",
		"push": "ok",
	}
)

type Request struct {
	Id     int    `json:"id"`
	Action string `json:"action"`
	Args   string `json:"args,omitempty"`
}

func (r *Request) String() string {
	s := "id:" + strconv.Itoa(r.Id) + ", action:" + r.Action
	if r.Args != "" {
		s += ", args: " + r.Args
	}
	return s
}

func parseRequest(message []byte, r *Request) error {
	if err := json.Unmarshal(message, r); err != nil {
		log.Println("Json unmarshal request params error")
		return err
	}
	return nil
}
