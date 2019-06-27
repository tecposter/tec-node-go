package post

import (
	"errors"
	"github.com/tecposter/tec-node-go/internal/draft"
	"github.com/tecposter/tec-node-go/internal/ws"
)

// commands
const (
	createCmd = "post.create"
	editCmd   = "post.edit"
	fetchCmd  = "post.fetch"
	commitCmd = "post.commit"
	listCmd   = "post.list"
)

// errors
var (
	ErrCmdNotFound  = errors.New("Command not found in post module")
	ErrDataDirEmpty = errors.New("dataDir cannot be empty")
	ErrUIDEmpty     = errors.New("uid cannot be empty")
)

// Service in post
type Service struct {
	dataDir string
	drftCtn *draft.Container
}

// NewService returns post.Service
func NewService(dataDir string) (*Service, error) {
	if dataDir == "" {
		return nil, ErrDataDirEmpty
	}

	drftCtn := draft.NewCtn(dataDir)
	svc := &Service{
		drftCtn: drftCtn}

	return svc, nil
}

// Close post service
func (svc *Service) Close() {
	svc.drftCtn.Close()
}

// HandleMsg response and request of wosocket message
func (svc *Service) HandleMsg(res *ws.Response, req *ws.Request) {
	switch req.Cmd() {
	case createCmd:
		svc.create(res, req)
	case editCmd:
		svc.edit(res, req)
	case fetchCmd:
		svc.fetch(res, req)
	case commitCmd:
		svc.commit(res, req)
	case listCmd:
		svc.list(res, req)
	default:
		res.Error(ErrCmdNotFound)
	}
}

/*
 *
 */
func (svc *Service) create(res *ws.Response, req *ws.Request) {
	uid := req.UID()
	if uid == "" {
		res.Error(ErrUIDEmpty)
		return
	}

	drfRepo, err := svc.drftCtn.Repo(uid)
	if err != nil {
		res.Error(err)
		return
	}

	pid, err := drfRepo.Reg()
	if err != nil {
		res.Error(err)
	} else {
		res.Set("pid", pid)
	}
}
func (svc *Service) edit(res *ws.Response, req *ws.Request) {
	res.Set("echo", "post.edit")
}
func (svc *Service) fetch(res *ws.Response, req *ws.Request) {
	res.Set("echo", "post.fetch")
}
func (svc *Service) commit(res *ws.Response, req *ws.Request) {
	res.Set("echo", "post.commit")
}
func (svc *Service) list(res *ws.Response, req *ws.Request) {
	res.Set("echo", "post.list")
}
