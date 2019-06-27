package draft

import (
	"errors"
	"github.com/tecposter/tec-node-go/internal/com/dto"
	"github.com/tecposter/tec-node-go/internal/com/store"
	"github.com/tecposter/tec-node-go/internal/com/uuid"
	"path"
)

var (
	ErrKeyNotFound = errors.New("Key not found")
)

var (
	errPIDEmpty = errors.New("PID cannot be empty")
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

func (repo *Repository) Reg() (string, error) {
	id, err := uuid.NewID()
	if err != nil {
		return "", err
	}
	drft := newDrft(id, "", "")
	err = repo.saveDrft(drft)
	if err != nil {
		return "", err
	}

	return drft.PID.Base58(), nil
}

func (repo *Repository) save(pidStr string, typ string, body string) error {
	if pidStr == "" {
		return errPIDEmpty
	}

	pid := dto.FromBase58(pidStr)
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
	drftData, err := drft.Marshal()
	if err != nil {
		return err
	}

	err = repo.db.Set(drft.PID.Bytes(), drftData)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) fetch(pidStr string) (*draft, error) {
	pid := dto.FromBase58(pidStr)
	res, err := repo.db.Get(pid.Bytes())
	if err != nil {
		return nil, err
	}

	var d draft
	err = d.Unmarshal(res)
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
		err := d.Unmarshal(val)
		if err != nil {
			return err
		}

		/*
			fmt.Println(d)
			keyID := dto.ID(key)
			fmt.Println(keyID.Base58())
		*/

		arr = append(arr, draftItem{
			PID:     d.PID,
			Changed: d.Changed,
			Title:   d.Title()})

		return nil
	})

	return arr, err
}

func (repo *Repository) delete(pidStr string) error {
	if pidStr == "" {
		return errPIDEmpty
	}
	pid := dto.FromBase58(pidStr)
	return repo.db.Delete(pid.Bytes())
}

func (repo *Repository) Close() {
	repo.db.Close()
}
