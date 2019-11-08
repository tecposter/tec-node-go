package ws

import (
	"encoding/json"
	"errors"
)

type responseDTO struct {
	CMD    string                 `json:"cmd"`
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

// Response websoket response
type Response struct {
	inner *responseDTO
	err   error
}

const (
	errorStatus = "error"
	okStatus    = "ok"
)

var (
	errEmpty = errors.New("Empty")
)

func newRes(cmd string) *Response {
	r := &Response{
		inner: &responseDTO{
			CMD:    cmd,
			Status: errorStatus,
			Data:   make(map[string]interface{}),
		},
		err: nil,
	}
	return r
}

func newErrRes(err error) *Response {
	r := newRes("unknown")
	r.SetErr(err)
	return r
}

// Marshal returns JSON encoding
func (res *Response) Marshal() ([]byte, error) {
	if res.err == nil {
		res.inner.Status = okStatus
	} else {
		res.inner.Status = errorStatus
		res.Set("error", res.err.Error())
	}
	b, err := json.Marshal(res.inner)
	return b, err
}

// SetErr sets error response
func (res *Response) SetErr(err error) {
	res.err = err
}

// Set sets response value of key
func (res *Response) Set(key string, val interface{}) {
	res.inner.Data[key] = val
}

// Has checks whether has key
func (res *Response) Has(key string) bool {
	_, ok := res.inner.Data[key]
	return ok
}

// Load loads response data
func (res *Response) Load(d map[string]interface{}) {
	for k, v := range d {
		res.Set(k, v)
	}
}
