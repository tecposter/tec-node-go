package draft

import (
	"errors"

	"github.com/tecposter/tec-node-go/src/ws"
)

const (
	cmdSave   = "draft.save"
	cmdFetch  = "draft.fetch"
	cmdList   = "draft.list"
	cmdDelete = "draft.delete"
	cmdHas    = "draft.has"
)

var (
	errCmdNotFound    = errors.New("Command not found in draft module")
	errRequirePostID  = errors.New("Require post id")
	errRequireContent = errors.New("Require content")
)

// Handle handle ws connection
func Handle(c *ws.Connection) {
	switch c.Req().CMD() {
	case cmdSave:
		save(c)
	case cmdFetch:
		fetch(c)
	case cmdList:
		list(c)
	case cmdDelete:
		delete(c)
	case cmdHas:
		has(c)
	default:
		c.Res().SetErr(errCmdNotFound)
	}
}

func save(c *ws.Connection) {
	req := c.Req()
	res := c.Res()

	postIDBase58, ok := req.Param("postID")
	if !ok {
		res.SetErr(errRequirePostID)
		return
	}
	content, ok := req.Param("content")
	if !ok {
		res.SetErr(errRequireContent)
		return
	}

	err := newServ(c).save(postIDBase58.(string), content.(string))
	if err != nil {
		res.SetErr(err)
	}
}

func fetch(c *ws.Connection) {
	req := c.Req()
	res := c.Res()

	postIDBase58, ok := req.Param("postID")
	if !ok {
		res.SetErr(errRequirePostID)
		return
	}

	d, err := newServ(c).fetch(postIDBase58.(string))
	if err != nil {
		res.SetErr(err)
		return
	}
	res.Set("draft", d)
}

func list(c *ws.Connection) {
	res := c.Res()

	list, err := newServ(c).list()
	if err != nil {
		res.SetErr(err)
		return
	}
	res.Set("drafts", list)
}

func delete(c *ws.Connection) {
	req := c.Req()
	res := c.Res()

	postIDBase58, ok := req.Param("postID")
	if !ok {
		res.SetErr(errRequirePostID)
		return
	}

	err := newServ(c).delete(postIDBase58.(string))
	if err != nil {
		res.SetErr(err)
		return
	}
}

func has(c *ws.Connection) {
	req := c.Req()
	res := c.Res()

	postIDBase58, ok := req.Param("postID")
	if !ok {
		res.SetErr(errRequirePostID)
		return
	}

	has, err := newServ(c).has(postIDBase58.(string))
	if err != nil {
		res.SetErr(err)
		return
	}
	res.Set("postID", postIDBase58)
	res.Set("has", has)
}
