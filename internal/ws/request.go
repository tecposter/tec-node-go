package ws

import (
	"encoding/json"
)

type Request struct {
	Cmd string `json:"cmd"`
	Token string `json:"token"`
	Params map[string]interface{} `json:"params"`
}

func ParseRequest(txt string) (*Request, error) {
		var req Request;
		if err := json.Unmarshal([]byte(txt), &req); err != nil {
			return nil, err
		}

		return &req, nil
}
