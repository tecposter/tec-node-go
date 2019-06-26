package user

import (
	//"log"
	"errors"
	"github.com/tecposter/tec-node-go/internal/com/uuid"
	"github.com/tecposter/tec-node-go/internal/ws"
)

const (
	regCmd    = "user.reg"
	loginCmd  = "user.login"
	logoutCmd = "user.logout"
)

const (
	tokenByteSize = 36
	lenMin        = 7
)

// errors
var (
	ErrUsernameEmpty    = errors.New("Usernaame cannot be empty")
	ErrUsernameTooShort = errors.New("Username too short - minimum length is 6")
	ErrUsernameExists   = errors.New("Username already exists")
	ErrPasswordTooShort = errors.New("Password too short - minimum length is 7")
	ErrPasswordNotMatch = errors.New("Password not match")
	ErrPasswordEmpty    = errors.New("password cannot be empty")
	ErrEmailExists      = errors.New("Email already exists")
	ErrEmailFormat      = errors.New("Error eamil format")
	ErrEmailNotFound    = errors.New("Email not found")
	ErrEmailEmpty       = errors.New("Email cannot be empty")
	ErrCmdNotFound      = errors.New("Command not found in user module")
	ErrNotLogin         = errors.New("Not Login")
)

type Service struct {
	repo *repository
}

func NewService(userDataDir string) (*Service, error) {
	repo, err := newRepo(userDataDir)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		repo: repo}

	return svc, nil
}

func (svc *Service) Close() {
	svc.repo.Close()
}

func (svc *Service) HandleMsg(res *ws.Response, req *ws.Request) {
	//log.Printf("ws.Request: %+v\n", req)
	switch req.Cmd() {
	case regCmd:
		svc.reg(res, req)
	case loginCmd:
		svc.login(res, req)
	case logoutCmd:
		svc.logout(res, req)
	default:
		res.Error(ErrCmdNotFound)
	}
}

/*
 * ---
 */

func (svc *Service) reg(res *ws.Response, req *ws.Request) {
	email := req.ParamStr("email")
	if svc.repo.hasEmail(email) {
		res.Error(ErrEmailExists)
		return
	}

	username := req.ParamStr("username")
	if username == "" {
		res.Error(ErrUsernameEmpty)
		return
	}
	if len(username) < 6 {
		res.Error(ErrUsernameTooShort)
		return
	}
	if svc.repo.hasUsername(username) {
		res.Error(ErrUsernameEmpty)
	}

	password := req.ParamStr("password")
	if len(password) < 7 {
		res.Error(ErrPasswordTooShort)
		return
	}

	passhash, err := hashPassword(password)
	if err != nil {
		res.Error(err)
		return
	}

	uid, err := uuid.NewBase58()
	if err != nil {
		res.Error(err)
		return
	}

	err = svc.repo.saveUser(uid, email, username, passhash)
	if err != nil {
		res.Error(err)
		return
	}

	res.Load(map[string]interface{}{
		"uid":      uid,
		"email":    email,
		"username": username})
}

func (svc *Service) login(res *ws.Response, req *ws.Request) {
	email := req.ParamStr("email")
	if email == "" {
		res.Error(ErrEmailEmpty)
		return
	}

	password := req.ParamStr("password")
	if password == "" {
		res.Error(ErrPasswordEmpty)
		return
	}

	uid := svc.repo.fetchUidByEmail(email)
	if uid == "" {
		res.Error(ErrEmailNotFound)
		return
	}

	usr, err := svc.repo.fetchUser(uid)
	if err != nil {
		res.Error(err)
		return
	}

	if !checkPasswordHash(password, usr.Passhash) {
		res.Error(ErrPasswordNotMatch)
		return
	}

	req.SetUid(uid)
}

func (svc *Service) logout(res *ws.Response, req *ws.Request) {
	req.RemoveUid()
}
