package ws

import (
	"golang.org/x/net/websocket"
)

type connHandler struct {
	conn *websocket.Conn
	bag map[string]interface{}
}

func newConnHandler(conn *websocket.Conn) *connHandler {
	return &connHandler{
		conn: conn,
		bag: make(map[string]interface{})}
}

func (hdl *connHandler) set(key string, val interface{}) {
	hdl.bag[key] = val
}

func (hdl *connHandler) get(key string) (interface{}, bool) {
	v, ok := hdl.bag[key]
	return v, ok
}

func (hdl *connHandler) remove(key string) {
	delete(hdl.bag, key)
}
