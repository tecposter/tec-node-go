package ws

import (
	"encoding/json"
)

type Response struct {
	Cmd string `json:"cmd"`
	Status string `json:"status"`
	Data map[string]interface{} `json:"data"`
}

func NewResponse(cmd string) *Response {
	return &Response{
		Status: "ok",
		Data: make(map[string]interface{}),
		Cmd: cmd}
}

func (r *Response) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Response) Error(msg string) {
	r.Status = "error"
	r.Data["message"] = msg
}
