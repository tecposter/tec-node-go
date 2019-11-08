package draft

import (
	"database/sql"

	"github.com/tecposter/tec-node-go/lib/dto"
	"github.com/tecposter/tec-node-go/src/ws"
)

type service struct {
	conn *ws.Connection
}

func newServ(c *ws.Connection) *service {
	return &service{conn: c}
}

func (serv *service) DB() *sql.DB {
	return serv.conn.DB()
}

func (serv *service) save(postIDBase58 string, content string) error {
	postID := dto.Base58ToID(postIDBase58)
	return newRepo(serv.DB()).save(postID, content)
}

func (serv *service) fetch(postIDBase58 string) (*draftDTO, error) {
	postID := dto.Base58ToID(postIDBase58)
	return newRepo(serv.DB()).fetch(postID)
}

func (serv *service) list() ([]draftDTO, error) {
	return newRepo(serv.DB()).list()
}

func (serv *service) delete(postIDBase58 string) error {
	postID := dto.Base58ToID(postIDBase58)
	return newRepo(serv.DB()).delete(postID)
}

func (serv *service) has(postIDBase58 string) (bool, error) {
	postID := dto.Base58ToID(postIDBase58)
	return newRepo(serv.DB()).has(postID)
}
