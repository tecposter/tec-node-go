package ws

import (
	"golang.org/x/net/websocket"
)

// Connection wraps websocket.Conn
type Connection struct {
	inner *websocket.Conn
	bag   map[string]interface{}
}

func newCollection(inner *websocket.Conn) *Connection {
	return &Connection{
		inner: inner,
		bag:   make(map[string]interface{})}
}

// Set adds a key-value pair to the bag of Connection
func (conn *Connection) Set(key string, val interface{}) {
	conn.bag[key] = val
}

// Get returns the value stored in the bag for a key
func (conn *Connection) Get(key string) (interface{}, bool) {
	val, ok := conn.bag[key]
	return val, ok
}

// Delete deletes the value for a key
func (conn *Connection) Delete(key string) {
	delete(conn.bag, key)
}
