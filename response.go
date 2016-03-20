package main

import (
	"encoding/json"
	"strconv"
	"time"
)

var (
	ServerMessage = map[string]string{
		"existed":         time.Now().Format(time.RFC822) + " you're forced to exit",
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
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
	Type    string `json:"type"`
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
		Code:    ServerMessageCode[msg],
		Message: ServerMessage[msg],
		Type:    rtype,
	}
}

func newActionResponse(action string) *Response {
	return &Response{
		Code:    0,
		Message: ReqResp[action],
		Type:    ResponseReply,
	}
}

func newPushResponse(msg string) *Response {
	return &Response{
		Code:    0,
		Message: msg,
		Type:    ResponsePush,
	}
}
