package app

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/tecposter/tec-node-go/app/draft"
	"github.com/tecposter/tec-node-go/app/post"
	"github.com/tecposter/tec-node-go/lib/db/sqlite3"
	"github.com/tecposter/tec-node-go/lib/ws"
)

// errors
var (
	ErrModuleNotFound = errors.New("Module not found")
)

const (
	postModule  = "post"
	draftModule = "draft"
)

var upgrader = websocket.Upgrader{}

type connection struct {
	inner     *websocket.Conn
	db        *sql.DB
	draftCtrl *draft.Controller
	postCtrl  *post.Controller
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
		inner:     inner,
		db:        db,
		postCtrl:  post.NewCtrl(db),
		draftCtrl: draft.NewCtrl(db),
		handleErr: func(err error) {
			log.Println("default handleErr: ", err)
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

		res, err := conn.dispatch(req)
		if err != nil {
			conn.handleErr(err)
			continue
		}

		b, err := res.Marshal()
		if err != nil {
			conn.handleErr(err)
			continue
		}

		err = conn.inner.WriteMessage(mt, b)
		if err != nil {
			conn.handleErr(err)
			continue
		}
	}
}

func (conn *connection) onErr(fn func(error)) {
	conn.handleErr = fn
}

func (conn *connection) dispatch(req ws.IRequest) (ws.IResponse, error) {
	switch req.Module() {
	case draftModule:
		return conn.draftCtrl.Handle(req)
	case postModule:
		return conn.postCtrl.Handle(req)
	default:
		return nil, ErrModuleNotFound
	}
}
