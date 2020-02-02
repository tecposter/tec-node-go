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
	p, err := newServ(c).create()
	if err == nil {
		c.Res().Set("post", p)
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
	serv := newServ(c)
	err := serv.edit(postIDBase58)
	if err != nil {
		res.SetErr(err)
		return
	}
	p, err := serv.fetch(postIDBase58)
	if err != nil {
		res.SetErr(err)
		return
	}
	res.Set("post", p)
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

	postIDBase58, ok := req.ParamStr("postID")
	if !ok {
		res.SetErr(errRequirePostID)
		return
	}
	contentType, ok := req.ParamStr("contentType")
	if !ok {
		res.SetErr(errRequireContentType)
		return
	}
	content, ok := req.ParamStr("content")
	if !ok {
		res.SetErr(errRequireContent)
		return
	}

	serv := newServ(c)
	err := serv.commit(
		postIDBase58,
		contentType,
		content,
	)
	if err != nil {
		res.SetErr(err)
		return
	}

	p, err := serv.fetch(postIDBase58)
	if err != nil {
		res.SetErr(err)
		return
	}
	res.Set("post", p)
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

	searchResult := newServ(c).search(query.(string))
	c.Res().Set("searchResult", searchResult)
}
