package ws

import (
	"golang.org/x/net/websocket"
)

type Connection struct {
	inner *websocket.Conn
	bag map[string]interface{}
}

func newCollection(inner *websocket.Conn) *Connection {
	return &Connection{
		inner: inner,
		bag: make(map[string]interface{})}
}

func (conn *Connection) set(key string, val interface{}) {
	conn.bag[key] = val
}

func (conn *Connection) get(key string) (interface{}, bool) {
	v, ok := conn.bag[key]
	return v, ok
}

func (conn *Connection) remove(key string) {
	delete(conn.bag, key)
}
