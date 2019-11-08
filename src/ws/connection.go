package ws

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

// Connection websocket connection
type Connection struct {
	inner *websocket.Conn
	db    *sql.DB
	res   *Response
	req   *Request
	mt    int
}

// NewConn returns websocket connection
func NewConn(w http.ResponseWriter, r *http.Request, db *sql.DB) (*Connection, error) {
	inner, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return &Connection{
		inner: inner,
		db:    db,
	}, nil
}

// Handle handle websocket connection
func (c *Connection) Handle(dispatch func()) {
	for {
		err := c.recv()
		if err != nil {
			log.Println("Connection.recv: ", err)
			continue
		}

		dispatch()

		err = c.send()
		if err != nil {
			log.Println("Connection.send: ", err)
			continue
		}
	}
}

// Close close the websocket connection
func (c *Connection) Close() {
	c.inner.Close()
	c.DB().Close()
}

// Res returns websocket response
func (c *Connection) Res() *Response {
	return c.res
}

// Req returns websocket request
func (c *Connection) Req() *Request {
	return c.req
}

// DB returns database connection
func (c *Connection) DB() *sql.DB {
	return c.db
}

func (c *Connection) recv() error {
	var message []byte
	var err error

	c.mt, message, err = c.inner.ReadMessage()
	if err != nil {
		return err
	}
	c.req, err = unmarshalWSReq(message)
	if err != nil {
		return err
	}
	c.res = newRes(c.req.CMD())

	return nil
}

/*
func (c *Connection) dispatch() {
	switch req.Module() {
	case draftModule:
		log.Println("draft module")
	case postModule:
		post.Handle(c)
	default:
		res.SetErr(ErrModuleNotFound)
	}
}
*/

func (c *Connection) send() error {
	b, err := c.Res().Marshal()
	if err != nil {
		return err
	}
	err = c.inner.WriteMessage(c.mt, b)
	return err
}
