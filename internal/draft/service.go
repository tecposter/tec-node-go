package draft

import (
	"errors"
	"github.com/tecposter/tec-node-go/internal/ws"
)

const (
	fetchCmd = "draft.fetch"
	saveCmd  = "draft.save"
	listCmd  = "draft.list"
)

const (
	cmdNotFoundErr  = "Command not found in draft module"
	dataDirEmptyErr = "dataDir cannot be empty"
	uidEmptyErr     = "uid cannot be empty"
)

var (
	errPidNotFound = errors.New("Pid not found")
	errUIDEmpty    = errors.New("UID cannot be empty")
)

type Service struct {
	ctn *Container
}

func NewService(dataDir string) (*Service, error) {
	if dataDir == "" {
		return nil, errors.New(dataDirEmptyErr)
	}

	return &Service{
		ctn: NewCtn(dataDir)}, nil
}

func (svc *Service) Close() {
}

func (svc *Service) HandleMsg(res *ws.Response, req *ws.Request) {
	switch req.Cmd() {
	case fetchCmd:
		svc.fetch(res, req)
	case saveCmd:
		svc.save(res, req)
	case listCmd:
		svc.list(res, req)
	default:
		res.Error(cmdNotFoundErr)
	}
}

func (svc *Service) fetch(res *ws.Response, req *ws.Request) {
	pid := req.ParamStr("pid")
	if pid == "" {
		res.Error(errPidNotFound.Error())
		return
	}

	repo, err := getRepo(svc, req)
	if err != nil {
		res.Error(err.Error())
		return
	}

	drft, err := repo.fetch(pid)
	if err != nil {
		res.Error(err.Error())
		return
	}
	res.Set("draft", drft)
}

func (svc *Service) save(res *ws.Response, req *ws.Request) {
	pid := req.ParamStr("pid")
	typ := req.ParamStr("typ")
	body := req.ParamStr("body")

	repo, err := getRepo(svc, req)
	if err != nil {
		res.Error(err.Error())
		return
	}
	err = repo.save(pid, typ, body)
	if err != nil {
		res.Error(err.Error())
		return
	}
}

func (svc *Service) list(res *ws.Response, req *ws.Request) {
}

// local func
func assertUID(req *ws.Request) string {
	uid := req.Uid()
	if uid == "" {
		panic(errUIDEmpty)
	}
	return uid
}

func getRepo(svc *Service, req *ws.Request) (*Repository, error) {
	uid := req.Uid()
	if uid == "" {
		return nil, errUIDEmpty
	}
	repo, err := svc.ctn.Repo(uid)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
