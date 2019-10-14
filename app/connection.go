package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/tecposter/tec-node-go/lib/db/sqlite3"
)

var upgrader = websocket.Upgrader{}

type connection struct {
	inner     *websocket.Conn
	db        *sql.DB
	handleErr func(error)
}

func newConn(dataDir string, w http.ResponseWriter, r *http.Request) (*connection, error) {
	db, err := sqlite3.Open(dataDir)
	if err != nil {
		return nil, err
	}

	inner, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return &connection{
		inner: inner,
		db:    db,
		handleErr: func(err error) {
			panic(err)
		},
	}, nil
}

func (conn *connection) close() {
	conn.inner.Close()
	conn.inner.Close()
}

func (conn *connection) run() {
	for {
		mt, message, err := conn.inner.ReadMessage()
		if err != nil {
			conn.handleErr(err)
			continue
		}

		log.Printf("recv: %s", message)
		req, err := unmarshalWSReq(message)
		if err != nil {
			conn.handleErr(err)
			continue
		}
		b, err := req.Marshal()
		if err != nil {
			log.Println("r.Marshal: ", err)
			continue
		}

		err = conn.inner.WriteMessage(mt, b)
		if err != nil {
			log.Println("write: ", err)
			continue
		}
	}
}

func (conn *connection) onErr(fn func(error)) {
	conn.handleErr = fn
}
