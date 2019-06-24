package ws

import (
	"sync"
	"golang.org/x/net/websocket"
)

type Connection struct {
	inner *websocket.Conn
	bag sync.Map
}

func newCollection(inner *websocket.Conn) *Connection {
	return &Connection{
		inner: inner}
}

func (conn *Connection) Set(key string, val interface{}) {
	conn.bag.Store(key, val)
}

func (conn *Connection) Get(key string) (interface{}, bool) {
	return conn.bag.Load(key)
}

func (conn *Connection) Remove(key string) {
	conn.bag.Delete(key)
}
