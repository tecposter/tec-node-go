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

const (
	cmdUnknown = "cmdUnknown"
)

var upgrader = websocket.Upgrader{}

type connection struct {
	inner     *websocket.Conn
	db        *sql.DB
	draftCtrl *draft.Controller
	postCtrl  *post.Controller
	// errHook   func(error)
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
	}, nil
}

func (conn *connection) close() {
	conn.inner.Close()
	conn.inner.Close()
}

func (conn *connection) recv() (int, ws.IRequest, error) {
	mt, message, err := conn.inner.ReadMessage()
	if err != nil {
		return mt, nil, err
	}

	req, err := unmarshalWSReq(message)
	if err != nil {
		return mt, nil, err
	}
	return mt, req, nil
}

func (conn *connection) run() {
	for {
		mt, req, err := conn.recv()
		if err != nil {
			log.Println("conn.recv: ", err)
			continue
		}

		res := newRes(req.CMD())
		conn.dispatch(res, req)

		err = conn.send(mt, res)
		if err != nil {
			log.Println("conn.send: ", err)
			continue
		}
		/*
			mt, message, err := conn.inner.ReadMessage()
			if err != nil {
				conn.handleErr(cmdUnknown, err)
				continue
			}

			log.Printf("recv: %s", message)
			req, err := unmarshalWSReq(message)
			if err != nil {
				conn.handleErr(cmdUnknown, err)
				continue
			}
		*/

		/*
			cmd := res.CMD()
			res, err := conn.dispatch(req)
			if err != nil {
				conn.handleErr(cmd, err)
				continue
			}

			b, err := res.Marshal()
			if err != nil {
				conn.handleErr(cmd, err)
				continue
			}

			err = conn.inner.WriteMessage(mt, b)
			if err != nil {
				conn.handleErr(cmd, err)
				continue
			}
		*/
	}
}

/*
func (conn *connection) handleErr(cmd string, err error) {
	res := newRes(cmd)
	res.SetErr(err)
}
*/

func (conn *connection) send(mt int, res ws.IResponse) error {
	b, err := res.Marshal()
	if err != nil {
		return err
	}
	err = conn.inner.WriteMessage(mt, b)
	return err
}

func (conn *connection) dispatch(res ws.IResponse, req ws.IRequest) {
	switch req.Module() {
	case draftModule:
		conn.draftCtrl.Handle(res, req)
	case postModule:
		conn.postCtrl.Handle(res, req)
	default:
		res.SetErr(ErrModuleNotFound)
	}
}
