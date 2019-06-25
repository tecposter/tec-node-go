package post

import (
	"github.com/tecposter/tec-server-go/internal/com/iotool"
	"github.com/tecposter/tec-server-go/internal/draft"
	"path"
	"sync"
)

type draftRepoCtn struct {
	dataDir string
	inner   sync.Map
}

func newDrfRepoCtn(dataDir string) *draftRepoCtn {
	return &draftRepoCtn{
		dataDir: dataDir}
}

func (ctn *draftRepoCtn) Repo(uid string) (*draft.Repository, error) {
	if val, ok := ctn.inner.Load(uid); ok {
		return val.(*draft.Repository), nil
	}

	iotool.MkdirIfNotExist(path.Join(ctn.dataDir, uid))
	newRepo, err := draft.NewRepo(ctn.dataDir, uid)
	if err != nil {
		return nil, err
	}
	ctn.inner.Store(uid, newRepo)
	return newRepo, nil
}
