package post

import (
	"database/sql"

	"github.com/tecposter/tec-node-go/lib/dto"
	"github.com/tecposter/tec-node-go/lib/uuid"
)

type service struct {
	repo *repository
}

func newServ(db *sql.DB) *service {
	return &service{
		repo: newRepo(db),
	}
}

func (serv *service) create() (dto.ID, error) {
	postID, err := uuid.NewID()
	if err != nil {
		return postID, err
	}

	err = serv.repo.create(postID)
	return postID, err
}
