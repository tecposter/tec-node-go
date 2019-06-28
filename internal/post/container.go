package post

import (
	"github.com/tecposter/tec-node-go/internal/com/dto"
	"github.com/tecposter/tec-node-go/internal/com/iotool"
	"path"
	"sync"
)

type container struct {
	dataDir string
}

var repos sync.Map
var once sync.Once

func newCtn(dataDir string) *container {
	return &container{dataDir: dataDir}
}

func (ctn *container) close() {
	once.Do(func() {
		repos.Range(func(_, r interface{}) bool {
			r.(*repository).close()
			return true
		})
	})
}

func (ctn *container) repo(uid dto.ID) (*repository, error) {
	key := uid.Base58()

	if val, ok := repos.Load(key); ok {
		return val.(*repository), nil
	}

	iotool.MkdirIfNotExist(path.Join(ctn.dataDir, key))
	repo, err := newRepo(ctn.dataDir, key)
	if err != nil {
		return nil, err
	}

	repos.Store(key, repo)
	return repo, nil
}
