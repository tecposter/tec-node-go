package draft

import (
	"path"
	"github.com/tecposter/tec-server-go/internal/com/store"
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
