package draft

import (
	"errors"
	"github.com/tecposter/tec-node-go/internal/com/dto"
	"github.com/tecposter/tec-node-go/internal/com/store"
	"github.com/tecposter/tec-node-go/internal/com/uuid"
	"path"
)

// errors
var (
	ErrKeyNotFound = errors.New("Key not found")
	ErrPIDEmpty    = errors.New("PID cannot be empty")
)

// A Repository wraps *store.DB
type Repository struct {
	db *store.DB
}

// NewRepo returns a Repository object
func NewRepo(dataDir string, uid string) (*Repository, error) {
	draftDir := path.Join(dataDir, uid, "draft")
	db, err := store.Open(draftDir)
	if err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

// Close closes Repository
func (repo *Repository) Close() {
	repo.db.Close()
}

// Reg returns PID of a new created draft
func (repo *Repository) Reg() (dto.ID, error) {
	id, err := uuid.NewID()
	if err != nil {
		return id, err
	}
	drft := newDrft(id, dto.TypText, "")
	err = repo.saveDrft(drft)
	if err != nil {
		return id, err
	}

	return drft.PID, nil
}

func (repo *Repository) save(pid dto.ID, typ dto.ContentType, body string) error {
	if pid == nil {
		return ErrPIDEmpty
	}

	//pid := dto.FromBase58(pidStr)
	ok, err := repo.db.Has(pid.Bytes())

	if err != nil {
		return err
	}
	if !ok {
		return ErrKeyNotFound
	}

	drft := newDrft(pid, typ, body)
	return repo.saveDrft(drft)
}

func (repo *Repository) saveDrft(drft *draft) error {
	id, data, err := drft.marshalPair()
	if err != nil {
		return err
	}

	err = repo.db.Set(id, data)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) fetch(pid dto.ID) (*draft, error) {
	//pid := dto.FromBase58(pidStr)
	res, err := repo.db.Get(pid.Bytes())
	if err != nil {
		return nil, err
	}

	var d draft
	err = d.unmarshalPair(pid.Bytes(), res)
	//err = json.Unmarshal(res, &d)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (repo *Repository) list() ([]draftItem, error) {
	var arr []draftItem

	err := repo.db.NewIter().ForEach(func(key, val []byte) error {
		var d draft
		err := d.unmarshalPair(key, val)
		if err != nil {
			return err
		}

		arr = append(arr, draftItem{
			PID:     d.PID,
			Changed: d.Changed,
			Title:   d.Title()})

		return nil
	})

	return arr, err
}

func (repo *Repository) delete(pid dto.ID) error {
	if pid == nil {
		return ErrPIDEmpty
	}
	//pid := dto.FromBase58(pidStr)
	return repo.db.Delete(pid.Bytes())
}
