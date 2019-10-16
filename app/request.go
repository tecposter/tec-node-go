package app

import (
	"encoding/json"
	"strings"
)

type requestDTO struct {
	CMD    string                 `json:"cmd"`
	Params map[string]interface{} `json:"params"`
}

type wsRequest struct {
	inner *requestDTO
}

func unmarshalWSReq(data []byte) (*wsRequest, error) {
	var inner requestDTO
	err := json.Unmarshal(data, &inner)
	if err != nil {
		return nil, err
	}

	return &wsRequest{
		inner: &inner,
	}, nil
}

func (r *wsRequest) Marshal() ([]byte, error) {
	b, err := json.Marshal(r.inner)
	return b, err
}

func (r *wsRequest) CMD() string {
	return r.inner.CMD
}

func (r *wsRequest) Module() string {
	cmd := r.CMD()
	dotIndex := strings.Index(cmd, ".")
	if dotIndex <= 0 {
		return ""
	}
	return cmd[0:dotIndex]
}

func (r *wsRequest) Param(key string) (interface{}, bool) {
	val, ok := r.inner.Params[key]
	return val, ok
}
