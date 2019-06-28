package draft

import (
	"github.com/tecposter/tec-node-go/internal/com/dto"
	"github.com/tecposter/tec-node-go/internal/com/iotool"
	"path"
	"sync"
)

// A Container contains users' repositories
type Container struct {
	dataDir string
}

var repos sync.Map
var once sync.Once

// NewCtn returns a new Container object
func NewCtn(dataDir string) *Container {
	return &Container{
		dataDir: dataDir}
}

// Close closes all the repositories in a Container
func (ctn *Container) Close() {
	once.Do(func() {
		repos.Range(func(_, r interface{}) bool {
			r.(*Repository).Close()
			return true
		})
	})
}

// Repo returns a Repository object for a uid
func (ctn *Container) Repo(uid dto.ID) (*Repository, error) {
	key := uid.Base58()

	if val, ok := repos.Load(key); ok {
		return val.(*Repository), nil
	}

	uidBase58 := uid.Base58()
	iotool.MkdirIfNotExist(path.Join(ctn.dataDir, uidBase58))
	repo, err := NewRepo(ctn.dataDir, uidBase58)
	if err != nil {
		return nil, err
	}

	repos.Store(key, repo)
	return repo, nil
}
