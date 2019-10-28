package post

import (
	"database/sql"
	"strings"

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

func (serv *service) commit(postIDBase58 string, contentTypeStr string, content string) error {
	postID := dto.Base58ToID(postIDBase58)
	cid := dto.GenContentID(content)
	ct := toCT(contentTypeStr)

	err := serv.repo.saveContent(cid, ct, content)
	if err != nil {
		return err
	}

	commitID, err := uuid.NewID()
	if err != nil {
		return err
	}

	err = serv.repo.commit(commitID, postID, cid)
	return err
}

func toCT(typeStr string) int {
	switch strings.ToLower(typeStr) {
	case "markdown":
		return typeMarkdown
	case "md":
		return typeMarkdown
	case "html":
		return typeHTML
	default:
		return typeText
	}
}
