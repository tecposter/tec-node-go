package post

import (
	"errors"
	"github.com/tecposter/tec-node-go/internal/com/dto"
	"github.com/tecposter/tec-node-go/internal/com/store"
	"path"
)

// errors
var (
	ErrKeyNotFound = errors.New("Key not found")
	ErrPIDEmpty    = errors.New("PID cannot be empty")
)

type repository struct {
	cmtDB *store.DB
	pstDB *store.DB
}

func newRepo(dataDir string, key string) (*repository, error) {
	postDir := path.Join(dataDir, key, "post")
	pstDB, err := store.Open(postDir)
	if err != nil {
		return nil, err
	}

	commitDir := path.Join(dataDir, key, "commit")
	cmtDB, err := store.Open(commitDir)

	return &repository{pstDB: pstDB, cmtDB: cmtDB}, nil
}

func (repo *repository) close() {
	repo.cmtDB.Close()
	repo.pstDB.Close()
}

func (repo *repository) fetch(pcid dto.ID) (*commit, error) {
	res, err := repo.cmtDB.Get(pcid)
	if err != nil {
		return nil, err
	}

	var c commit
	err = c.unmarshalPair(pcid, res)
	//err = json.Unmarshal(res, &d)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (repo *repository) hasPID(pid dto.ID) (bool, error) {
	return repo.pstDB.Has(pid)
}

func (repo *repository) commit(pid dto.ID, typ string, body string) (*commit, error) {
	if pid == nil {
		return nil, ErrPIDEmpty
	}

	pst, cmt := newPostCommit(pid, dto.MakeContent(typ, body))

	pid, pstData, err := pst.marshalPair()
	if err != nil {
		return nil, err
	}
	pcid, cmtData, err := cmt.marshalPair()
	if err != nil {
		return nil, err
	}

	err = repo.pstDB.Set(pid, pstData)
	if err != nil {
		return nil, err
	}
	err = repo.cmtDB.Set(pcid, cmtData)
	if err != nil {
		return nil, err
	}

	return cmt, nil
}

func (repo *repository) list() ([]post, error) {
	var arr []post

	err := repo.pstDB.NewIter().ForEach(func(key, val []byte) error {
		var pst post
		err := pst.unmarshalPair(key, val)
		if err != nil {
			return err
		}
		arr = append(arr, pst)
		return nil
	})
	return arr, err
}
