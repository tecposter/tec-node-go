package post

import (
	"errors"
	"github.com/tecposter/tec-server-go/internal/com/uuid"
	"github.com/tecposter/tec-server-go/internal/ws"
)

const (
	createCmd = "post.create"
	editCmd   = "post.edit"
	fetchCmd  = "post.fetch"
	commitCmd = "post.commit"
	listCmd   = "post.list"

	cmdNotFoundErr  = "Command not found in post module"
	dataDirEmptyErr = "dataDir cannot be empty"
	uidEmptyErr     = "uid cannot be empty"
)

type Service struct {
	dataDir    string
	drfRepoCtn *draftRepoCtn
}

func NewService(dataDir string) (*Service, error) {
	if dataDir == "" {
		return nil, errors.New(dataDirEmptyErr)
	}

	drfRepoCtn := newDrfRepoCtn(dataDir)
	svc := &Service{
		drfRepoCtn: drfRepoCtn}

	return svc, nil
}

func (svc *Service) Close() {
}

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
		res.Error(cmdNotFoundErr)
	}
}

/*
 *
 */
func (svc *Service) create(res *ws.Response, req *ws.Request) {
	uid := req.Uid()
	if uid == "" {
		res.Error(uidEmptyErr)
		return
	}

	drfRepo, err := svc.drfRepoCtn.Repo(uid)
	if err != nil {
		res.Error(err.Error())
		return
	}

	pid, err := uuid.NewBase58()
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = drfRepo.Reg(pid)
	if err != nil {
		res.Error(err.Error())
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
