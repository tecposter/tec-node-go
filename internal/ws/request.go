package ws

import (
	"encoding/json"
)

type Request struct {
	Conn *Connection
	data *requestData
}

type requestData struct {
	Cmd    string                 `json:"cmd"`
	Token  string                 `json:"token"`
	Params map[string]interface{} `json:"params"`
}

func NewRequest(conn *Connection, txt string) (*Request, error) {
	var data requestData
	if err := json.Unmarshal([]byte(txt), &data); err != nil {
		return nil, err
	}

	return &Request{Conn: conn, data: &data}, nil
}

func (req *Request) Cmd() string {
	return req.data.Cmd
}

func (req *Request) ParamStr(key string) string {
	val, ok := req.data.Params[key]
	if !ok {
		return ""
	}

	return val.(string)
}

const (
	uidKey = "uid"
)

func (req *Request) SetUID(uid string) {
	req.Conn.Set(uidKey, uid)
}

func (req *Request) UID() string {
	if v, ok := req.Conn.Get(uidKey); ok {
		return v.(string)
	}
	return ""
}

func (req *Request) RemoveUID() {
	req.Conn.Remove(uidKey)
}
