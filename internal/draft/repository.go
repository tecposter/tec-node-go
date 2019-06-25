package draft

import (
	"github.com/tecposter/tec-node-go/internal/com/store"
	"path"
)

type Repository struct {
	db *store.DB
}

func NewRepo(dataDir string, uid string) (*Repository, error) {
	draftDir := path.Join(dataDir, uid, "draft")
	db, err := store.Open(draftDir)
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (repo *Repository) Reg(pid string) error {
	emptyItem := []byte("{}")
	return repo.db.Set([]byte(pid), emptyItem)
}
