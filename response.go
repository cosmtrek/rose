package main

import (
	"encoding/json"
	"strconv"
	"time"
)

var (
	ServerMessage = map[string]string{
		"existed":         "you're forced to be offline",
		"unknown_actions": "unknown actions",
	}
	ServerMessageCode = map[string]int{
		"existed":         1,
		"unknown_actions": 2,
	}
)

const (
	ResponseReply = "reply"
	ResponsePush  = "push"
)

type Response struct {
	Type      string `json:"type"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Error     string `json:"error,omitempty"`
	CreatedAt string `json:"created_at"`
}

func (r *Response) String() string {
	s := "code:" + strconv.Itoa(r.Code) + ", message:" + r.Message
	if r.Error != "" {
		s += ", error: " + r.Error
	}
	return s
}

func (r *Response) Json() string {
	b, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(b)
}

func newResponse(msg string, rtype string) *Response {
	return &Response{
		Type:      rtype,
		Code:      ServerMessageCode[msg],
		Message:   ServerMessage[msg],
		CreatedAt: time.Now().Format(time.RFC3339),
	}
}

func newActionResponse(action string) *Response {
	return &Response{
		Type:      ResponseReply,
		Code:      0,
		Message:   ReqResp[action],
		CreatedAt: time.Now().Format(time.RFC3339),
	}
}

func newPushResponse(msg string) *Response {
	return &Response{
		Type:      ResponsePush,
		Code:      0,
		Message:   msg,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
}
