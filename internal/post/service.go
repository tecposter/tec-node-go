package post

import (
	"errors"
	"github.com/tecposter/tec-node-go/internal/com/dto"
	"github.com/tecposter/tec-node-go/internal/com/uuid"
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
	ErrPIDNotExists = errors.New("pid not exists")
	ErrPIDNotFound  = errors.New("pid not found")
)

// Service in post
type Service struct {
	dataDir string
	drftCtn *draft.Container
	ctn     *container
}

// NewService returns post.Service
func NewService(dataDir string) (*Service, error) {
	if dataDir == "" {
		return nil, ErrDataDirEmpty
	}

	drftCtn := draft.NewCtn(dataDir)
	svc := &Service{drftCtn: drftCtn, ctn: newCtn(dataDir)}

	return svc, nil
}

// Close post service
func (svc *Service) Close() {
	svc.drftCtn.Close()
	svc.ctn.close()
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
	pid, err := uuid.NewID()
	if err != nil {
		res.Error(err)
		return
	}
	svc.reg(res, req, pid)
}
func (svc *Service) edit(res *ws.Response, req *ws.Request) {
	pid := dto.Base58ToID(req.ParamStr("pid"))
	uid := req.UID()
	if uid == nil {
		res.Error(ErrUIDEmpty)
		return
	}

	repo, err := svc.ctn.repo(uid)
	if err != nil {
		res.Error(err)
		return
	}
	ok, err := repo.hasPID(pid)
	if err != nil {
		res.Error(err)
		return
	}
	if !ok {
		res.Error(ErrPIDNotFound)
		return
	}

	svc.reg(res, req, pid)
}

func (svc *Service) reg(res *ws.Response, req *ws.Request, pid dto.ID) {
	uid := req.UID()
	if uid == nil {
		res.Error(ErrUIDEmpty)
		return
	}

	drfRepo, err := svc.drftCtn.Repo(uid)
	if err != nil {
		res.Error(err)
		return
	}
	err = drfRepo.Reg(pid)
	if err != nil {
		res.Error(err)
	} else {
		res.Set("pid", pid)
	}
}

func (svc *Service) fetch(res *ws.Response, req *ws.Request) {
	uid := req.UID()
	if uid == nil {
		res.Error(ErrUIDEmpty)
		return
	}

	repo, err := svc.ctn.repo(uid)
	if err != nil {
		res.Error(err)
		return
	}
	pcid := dto.Base58ToID(req.ParamStr("pcid"))
	cmt, err := repo.fetch(pcid)
	if err != nil {
		res.Error(err)
		return
	}

	res.Set("post", cmt)
}

func (svc *Service) commit(res *ws.Response, req *ws.Request) {
	uid := req.UID()
	if uid == nil {
		res.Error(ErrUIDEmpty)
		return
	}
	drftRepo, err := svc.drftCtn.Repo(uid)
	if err != nil {
		res.Error(err)
		return
	}

	pid := dto.Base58ToID(req.ParamStr("pid"))
	ok, err := drftRepo.Has(pid)
	if err != nil {
		res.Error(err)
		return
	}
	if !ok {
		res.Error(ErrPIDNotExists)
		return
	}

	repo, err := svc.ctn.repo(uid)
	if err != nil {
		res.Error(err)
		return
	}

	typ := req.ParamStr("typ")
	body := req.ParamStr("body")
	cmt, err := repo.commit(pid, typ, body)
	if err != nil {
		res.Error(err)
		return
	}

	err = drftRepo.Delete(pid)
	if err != nil {
		res.Error(err)
		return
	}

	res.Set("pcid", cmt.PCID)
}

func (svc *Service) list(res *ws.Response, req *ws.Request) {
	uid := req.UID()
	if uid == nil {
		res.Error(ErrUIDEmpty)
		return
	}
	repo, err := svc.ctn.repo(uid)
	if err != nil {
		res.Error(err)
		return
	}
	list, err := repo.list()
	if err != nil {
		res.Error(err)
		return
	}

	res.Set("list", list)
}
