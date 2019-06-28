package draft

import (
	"errors"
	"github.com/tecposter/tec-node-go/internal/com/dto"
	"github.com/tecposter/tec-node-go/internal/ws"
)

const (
	fetchCmd  = "draft.fetch"
	saveCmd   = "draft.save"
	listCmd   = "draft.list"
	deleteCmd = "draft.delete"
)

const (
	cmdNotFoundErr  = "Command not found in draft module"
	dataDirEmptyErr = "dataDir cannot be empty"
	uidEmptyErr     = "UID cannot be empty"
)

// errors in the draft service
var (
	ErrCmdNotFound  = errors.New("Command not found in draft module")
	ErrDataDirEmpty = errors.New("dataDir cannot be empty")
	ErrPidNotFound  = errors.New("Pid not found")
	ErrUIDEmpty     = errors.New("UID cannot be empty")
)

// A Service provides functions to handle messages from clients
type Service struct {
	ctn *Container
}

// NewService returns a Service
func NewService(dataDir string) (*Service, error) {
	if dataDir == "" {
		return nil, errors.New(dataDirEmptyErr)
	}

	return &Service{
		ctn: NewCtn(dataDir)}, nil
}

// Close closes a Service
func (svc *Service) Close() {
	svc.ctn.Close()
}

// HandleMsg handles messages from client
func (svc *Service) HandleMsg(res *ws.Response, req *ws.Request) {
	switch req.Cmd() {
	case fetchCmd:
		svc.fetch(res, req)
	case saveCmd:
		svc.save(res, req)
	case listCmd:
		svc.list(res, req)
	case deleteCmd:
		svc.delete(res, req)
	default:
		res.Error(ErrCmdNotFound)
	}
}

func (svc *Service) fetch(res *ws.Response, req *ws.Request) {
	pid := dto.Base58ToID(req.ParamStr("pid"))
	if pid == nil {
		res.Error(ErrPidNotFound)
		return
	}

	repo, err := getRepo(svc, req)
	if err != nil {
		res.Error(err)
		return
	}

	drft, err := repo.fetch(pid)
	if err != nil {
		res.Error(err)
		return
	}
	res.Set("draft", drft)
}

func (svc *Service) save(res *ws.Response, req *ws.Request) {
	pid := dto.Base58ToID(req.ParamStr("pid"))
	if pid == nil {
		res.Error(ErrPidNotFound)
		return
	}

	typ := dto.ParseContentType(req.ParamStr("typ"))
	body := req.ParamStr("body")

	repo, err := getRepo(svc, req)
	if err != nil {
		res.Error(err)
		return
	}
	err = repo.save(pid, typ, body)
	if err != nil {
		res.Error(err)
		return
	}
}

func (svc *Service) list(res *ws.Response, req *ws.Request) {
	repo, err := getRepo(svc, req)
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

func (svc *Service) delete(res *ws.Response, req *ws.Request) {
	repo, err := getRepo(svc, req)
	if err != nil {
		res.Error(err)
		return
	}

	pid := dto.Base58ToID(req.ParamStr("pid"))
	if pid == nil {
		res.Error(ErrPidNotFound)
		return
	}
	err = repo.delete(pid)
	if err != nil {
		res.Error(err)
		return
	}
}

// local func
func assertUID(req *ws.Request) dto.ID {
	uid := req.UID()
	if uid == nil {
		panic(ErrUIDEmpty)
	}
	return uid
}

func getRepo(svc *Service, req *ws.Request) (*Repository, error) {
	uid := req.UID()
	if uid == nil {
		return nil, ErrUIDEmpty
	}
	repo, err := svc.ctn.Repo(uid)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
