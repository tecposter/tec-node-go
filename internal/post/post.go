package post

import (
	"bytes"
	"fmt"
	"github.com/tecposter/tec-node-go/internal/com/bin"
	"github.com/tecposter/tec-node-go/internal/com/dto"
	"time"
)

const (
	hyphen      = '-'
	titleLenMax = 100
)

const (
	idSize   = 16
	pcidSize = 25 // idSize + hyphenSize + timeSize
	timeSize = 8
)

type commit struct {
	PCID    dto.ID      `json:"pcid"`
	PID     dto.ID      `json:"pid"`
	Created time.Time   `json:"created"`
	Cont    dto.Content `json:"cont"`
}

type post struct {
	PID     dto.ID    `json:"pid"`
	PCID    dto.ID    `json:"pcid"`
	Changed time.Time `json:"changed"`
	Title   string    `json:"title"`
}

func newPostCommit(pid dto.ID, cont dto.Content) (*post, *commit) {
	nw := time.Now()

	pref := append(pid, hyphen)
	pcid := dto.ID(append(pref, bin.TimeToBytes(nw)...))

	fmt.Println(pid)
	fmt.Println(pref)
	fmt.Println(pcid)

	cmt := commit{
		PCID:    pcid,
		PID:     pid,
		Created: nw,
		Cont:    cont}

	pst := post{
		PID:     pid,
		PCID:    pcid,
		Changed: nw,
		Title:   cont.Title(titleLenMax)}

	return &pst, &cmt
}

func (c *commit) marshalPair() ([]byte, []byte, error) {
	pcid := c.PCID
	pid := c.PID
	created := bin.TimeToBytes(c.Created)
	typ := byte(c.Cont.Typ)
	body := []byte(c.Cont.Body)

	var buf bytes.Buffer
	buf.Write(pid)
	buf.Write(created)
	buf.WriteByte(typ)
	buf.Write(body)

	return pcid, buf.Bytes(), nil
}

func (c *commit) unmarshalPair(pcid, data []byte) error {
	pid := data[:idSize]
	created := bin.BytesToTime(data[idSize : idSize+timeSize])
	typ := dto.ContentType(data[idSize+timeSize])
	body := string(data[idSize+timeSize+1:])

	c.PCID = pcid
	c.PID = pid
	c.Created = created
	c.Cont.Typ = typ
	c.Cont.Body = body

	return nil
}

func (p *post) marshalPair() ([]byte, []byte, error) {
	pid := p.PID.Bytes()
	pcid := p.PCID.Bytes()
	changed := bin.TimeToBytes(p.Changed)
	title := []byte(p.Title)

	return pid, append(pcid, append(changed, title...)...), nil
}

func (p *post) unmarshalPair(pid, data []byte) error {
	pcid := data[:pcidSize]
	changed := bin.BytesToTime(data[pcidSize : pcidSize+timeSize])
	title := string(data[pcidSize+timeSize:])

	p.PID = pid
	p.PCID = pcid
	p.Changed = changed
	p.Title = title

	return nil
}
