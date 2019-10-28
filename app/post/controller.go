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
	errCmdNotFound        = errors.New("Command not found in post module")
	errRequirePostID      = errors.New("Require post id")
	errRequireContent     = errors.New("Require content")
	errRequireContentType = errors.New("Require content type")
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
	case cmdCommit:
		ctrl.commit(res, req)
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

func (ctrl *Controller) commit(res ws.IResponse, req ws.IRequest) {
	postIDBase58, ok := req.Param("postID")
	if !ok {
		res.SetErr(errRequirePostID)
		return
	}
	contentType, ok := req.Param("contentType")
	if !ok {
		res.SetErr(errRequireContentType)
		return
	}
	content, ok := req.Param("content")
	if !ok {
		res.SetErr(errRequireContent)
		return
	}

	err := ctrl.serv.commit(
		postIDBase58.(string),
		contentType.(string),
		content.(string),
	)
	if err != nil {
		res.SetErr(err)
	}
}
