package app

import (
	"encoding/json"
)

// IRequest request interface
type IRequest interface {
	CMD() string
	Param(string) (interface{}, bool)
}

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

func (r *wsRequest) Param(key string) (interface{}, bool) {
	val, ok := r.inner.Params[key]
	return val, ok
}
