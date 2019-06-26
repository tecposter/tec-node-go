package draft

import (
	"github.com/tecposter/tec-node-go/internal/com/iotool"
	"path"
	"sync"
)

type Container struct {
	dataDir string
}

var repos sync.Map
var once sync.Once

func NewCtn(dataDir string) *Container {
	return &Container{
		dataDir: dataDir}
}

func (ctn *Container) Close() {
	once.Do(func() {
		repos.Range(func(_, r interface{}) bool {
			r.(*Repository).Close()
			return true
		})
	})
}

func (ctn *Container) Repo(uid string) (*Repository, error) {
	if val, ok := repos.Load(uid); ok {
		return val.(*Repository), nil
	}

	iotool.MkdirIfNotExist(path.Join(ctn.dataDir, uid))
	repo, err := NewRepo(ctn.dataDir, uid)
	if err != nil {
		return nil, err
	}

	repos.Store(uid, repo)
	return repo, nil
}
