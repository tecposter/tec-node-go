package ws

import (
	"encoding/json"
	"strings"
)

type requestDTO struct {
	CMD    string                 `json:"cmd"`
	Params map[string]interface{} `json:"params"`
}

// Request ws request
type Request struct {
	inner *requestDTO
}

func unmarshalWSReq(data []byte) (*Request, error) {
	var inner requestDTO
	err := json.Unmarshal(data, &inner)
	if err != nil {
		return nil, err
	}

	return &Request{
		inner: &inner,
	}, nil
}

// Marshal returns JSON encoding
func (r *Request) Marshal() ([]byte, error) {
	b, err := json.Marshal(r.inner)
	return b, err
}

// CMD returns request command
func (r *Request) CMD() string {
	return r.inner.CMD
}

// Module return request module
func (r *Request) Module() string {
	cmd := r.CMD()
	dotIndex := strings.Index(cmd, ".")
	if dotIndex <= 0 {
		return ""
	}
	return cmd[0:dotIndex]
}

// Param returns request parameter of key
func (r *Request) Param(key string) (interface{}, bool) {
	val, ok := r.inner.Params[key]
	return val, ok
}

// ParamStr returns parameter as string by key
func (r *Request) ParamStr(key string) (string, bool) {
	val, ok := r.inner.Params[key]
	return strings.TrimSpace(val.(string)), ok
}
