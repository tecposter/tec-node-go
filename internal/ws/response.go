package ws

import (
	"encoding/json"
)

type Response struct {
	Cmd    string                 `json:"cmd"`
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

const (
	errorStatus = "error"
)

func NewResponse(cmd string) *Response {
	return &Response{
		Status: "ok",
		Data:   make(map[string]interface{}),
		Cmd:    cmd}
}

func (r *Response) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Response) Error(err error) {
	r.Status = errorStatus
	r.Data["message"] = err.Error()
}

func (res *Response) Set(key string, val interface{}) {
	res.Data[key] = val
}

func (res *Response) Load(d map[string]interface{}) {
	for k, v := range d {
		res.Set(k, v)
	}
}

func (res *Response) Has(key string) bool {
	_, ok := res.Data[key]
	return ok
}
