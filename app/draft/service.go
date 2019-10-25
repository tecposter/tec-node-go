package draft

import (
	"database/sql"

	"github.com/tecposter/tec-node-go/lib/dto"
)

type service struct {
	repo *repository
}

func newServ(db *sql.DB) *service {
	return &service{repo: newRepo(db)}
}

func (serv *service) save(postIDBase58 string, content string) error {
	postID := dto.Base58ToID(postIDBase58)
	return serv.repo.save(postID, content)
}

func (serv *service) fetch(postIDBase58 string) (*draftDTO, error) {
	postID := dto.Base58ToID(postIDBase58)
	return serv.repo.fetch(postID)
}

func (serv *service) list() ([]draftDTO, error) {
	return serv.repo.list()
}
