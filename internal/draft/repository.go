package draft

import (
	"encoding/json"
	"errors"
	"github.com/tecposter/tec-node-go/internal/com/store"
	"log"
	"path"
	"time"
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

func (repo *Repository) Reg(pid string) error {
	if pid == "" {
		return errPIDEmpty
	}
	d := draft{
		PID:     pid,
		Changed: time.Now(),
		Cont: content{
			Typ:  "",
			Body: ""}}
	bs, err := json.Marshal(d)
	if err != nil {
		return err
	}

	return repo.db.Set([]byte(pid), bs)
}

func (repo *Repository) save(pid string, typ string, body string) error {
	if pid == "" {
		return errPIDEmpty
	}

	key := []byte(pid)
	ok, err := repo.db.Has(key)
	if err != nil {
		return err
	}
	if !ok {
		return ErrKeyNotFound
	}

	log.Println(pid, typ, body)

	cont := content{
		Typ:  typ,
		Body: body}
	drft := draft{
		PID:     pid,
		Changed: time.Now(),
		Cont:    cont}

	drftData, err := json.Marshal(drft)
	if err != nil {
		return err
	}

	err = repo.db.Set(key, drftData)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) fetch(pid string) (*draft, error) {
	res, err := repo.db.Get([]byte(pid))
	if err != nil {
		return nil, err
	}

	var d draft
	err = json.Unmarshal(res, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (repo *Repository) list() ([]draftItem, error) {
	var arr []draftItem

	err := repo.db.NewIter().ForEach(func(key, val []byte) error {
		var d draft
		err := json.Unmarshal(val, &d)
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

func (repo *Repository) delete(pid string) error {
	if pid == "" {
		return errPIDEmpty
	}
	return repo.db.Delete([]byte(pid))
}

func (repo *Repository) Close() {
	repo.db.Close()
}
