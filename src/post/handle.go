package post

import (
	"errors"

	"github.com/tecposter/tec-node-go/src/ws"
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

// Handle handle websoket response and request
func Handle(c *ws.Connection) {
	switch c.Req().CMD() {
	case cmdCreate:
		create(c)
	case cmdEdit:
		edit(c)
	case cmdFetch:
		fetch(c)
	case cmdCommit:
		commit(c)
	case cmdList:
		list(c)
	case cmdSearch:
		search(c)
	default:
		c.Res().SetErr(errCmdNotFound)
	}
}

func create(c *ws.Connection) {
	postID, err := newServ(c).create()
	if err == nil {
		c.Res().Set("postID", postID.Base58())
	} else {
		c.Res().SetErr(err)
	}
}

func edit(c *ws.Connection) {
	req := c.Req()
	res := c.Res()

	postIDBase58, ok := req.ParamStr("postID")
	if !ok {
		res.SetErr(errRequirePostID)
		return
	}
	err := newServ(c).edit(postIDBase58)
	if err != nil {
		res.SetErr(err)
		return
	}
	res.Set("postID", postIDBase58)
}

func fetch(c *ws.Connection) {
	req := c.Req()
	res := c.Res()

	postIDBase58, ok := req.ParamStr("postID")
	if !ok {
		res.SetErr(errRequirePostID)
		return
	}
	p, err := newServ(c).fetch(postIDBase58)
	if err != nil {
		res.SetErr(err)
		return
	}
	res.Set("post", p)
}

func commit(c *ws.Connection) {
	req := c.Req()
	res := c.Res()

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

	err := newServ(c).commit(
		postIDBase58.(string),
		contentType.(string),
		content.(string),
	)
	if err != nil {
		res.SetErr(err)
	}
}

func list(c *ws.Connection) {
	res := c.Res()

	list, err := newServ(c).list()
	if err != nil {
		res.SetErr(err)
		return
	}
	res.Set("posts", list)
}

func search(c *ws.Connection) {
	query, ok := c.Req().Param("query")
	if !ok {
		// res.SetErr(errRequireQuery)
		// return
		query = ""
	}

	newServ(c).search(query.(string))
}
