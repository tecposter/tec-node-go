package app

import (
	"encoding/json"
	"errors"
)

type responseDTO struct {
	CMD    string                 `json:"cmd"`
	Status string                 `json:"status"`
	Error  string                 `json:"error"`
	Data   map[string]interface{} `json:"data"`
}

type wsResponse struct {
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

func newRes(cmd string) *wsResponse {
	r := &wsResponse{
		inner: &responseDTO{
			CMD:    cmd,
			Status: errorStatus,
			Error:  "",
			Data:   make(map[string]interface{}),
		},
		err: nil,
	}
	return r
}

func newErrRes(err error) *wsResponse {
	r := newRes("unknown")
	r.SetErr(err)
	return r
}

func (res *wsResponse) Marshal() ([]byte, error) {
	if res.err == nil {
		res.inner.Status = okStatus
		res.inner.Error = ""
	} else {
		res.inner.Status = errorStatus
		res.inner.Error = res.err.Error()
	}
	b, err := json.Marshal(res.inner)
	return b, err
}

func (res *wsResponse) SetErr(err error) {
	res.err = err
}

func (res *wsResponse) Set(key string, val interface{}) {
	res.inner.Data[key] = val
}

func (res *wsResponse) Has(key string) bool {
	_, ok := res.inner.Data[key]
	return ok
}

func (res *wsResponse) Load(d map[string]interface{}) {
	for k, v := range d {
		res.Set(k, v)
	}
}
