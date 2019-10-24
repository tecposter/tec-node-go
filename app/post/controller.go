package post

import (
	"database/sql"
	"errors"
	// "log"

	"github.com/tecposter/tec-node-go/lib/ws"
)

const (
	cmdCreate = "post.create"
	cmdEdit   = "post.edit"
	cmdFetch  = "post.fetch"
	cmdCommit = "post.commit"
	cmdList   = "post.list"
	cmdSearch = "post.search"
)

var (
	errCmdNotFound = errors.New("Command not found in post module")
)

// Controller in post
type Controller struct {
	serv *service
	// db *sql.DB
}

// NewCtrl return Ctroller instance
func NewCtrl(db *sql.DB) *Controller {
	return &Controller{serv: newServ(db)}
}

// Handle handle ws response and request in post
func (ctrl *Controller) Handle(res ws.IResponse, req ws.IRequest) {
	switch req.CMD() {
	case cmdCreate:
		ctrl.create(res)
	default:
		res.SetErr(errCmdNotFound)
	}
}

func (ctrl *Controller) create(res ws.IResponse) {
	postID, err := ctrl.serv.create()
	if err != nil {
		res.SetErr(err)
		return
	}

	res.Set("postID", postID.Base58())
}
