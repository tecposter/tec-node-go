package post

import (
	"database/sql"
	"strings"

	"github.com/tecposter/tec-node-go/lib/dto"
	"github.com/tecposter/tec-node-go/lib/uuid"
	"github.com/tecposter/tec-node-go/src/searcher"
	"github.com/tecposter/tec-node-go/src/ws"
)

type service struct {
	conn *ws.Connection
}

func newServ(c *ws.Connection) *service {
	return &service{
		conn: c,
	}
}

func (s *service) DB() *sql.DB {
	return s.conn.DB()
}

func (s *service) create() (*postDTO, error) {
	postID, err := uuid.NewID()
	if err != nil {
		return nil, err
	}

	post, err := newRepo(s.DB()).create(postID)
	return post, err
}

func (s *service) edit(postIDBase58 string) error {
	postID := dto.Base58ToID(postIDBase58)
	err := newRepo(s.DB()).edit(postID)
	return err
}

func (s *service) fetch(postIDBase58 string) (*postDTO, error) {
	postID := dto.Base58ToID(postIDBase58)
	return newRepo(s.DB()).fetch(postID)
}

func (s *service) commit(postIDBase58 string, contentType string, content string) error {
	postID := dto.Base58ToID(postIDBase58)
	contentID := dto.GenContentID(content)
	contentTypeID := toContentTypeID(contentType)

	repo := newRepo(s.DB())
	err := repo.saveContent(contentID, contentTypeID, content)
	if err != nil {
		return err
	}

	commitID, err := uuid.NewID()
	if err != nil {
		return err
	}

	err = repo.commit(commitID, postID, contentID)
	if err != nil {
		return err
	}

	searcher.Index(postID.Base58(), content)
	return nil
}

func (s *service) list() ([]postItemDTO, error) {
	return newRepo(s.DB()).list()
}

func (s *service) search(query string) []searcher.Doc {
	rs := searcher.Search(query)
	return rs
}

func toContentTypeID(contentType string) int {
	switch strings.ToLower(contentType) {
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
