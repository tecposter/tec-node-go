package user

import (
	//"log"
	"github.com/tecposter/tec-server-go/internal/com/uuid"
	"github.com/tecposter/tec-server-go/internal/ws"
)

const (
	regCmd    = "user.reg"
	loginCmd  = "user.login"
	logoutCmd = "user.logout"

	usernameEmptyErr    = "Usernaame cannot be empty"
	usernameTooShortErr = "Username too short - minimum length is 6"
	usernameExistsErr   = "Username already exists"

	passwordTooShortErr = "Password too short - minimum length is 7"
	passwordNotMatchErr = "Password not match"
	passwordEmptyErr    = "password cannot be empty"

	emailExistsErr   = "Email already exists"
	emailFormatErr   = "Error eamil format"
	emailNotFoundErr = "Email not found"
	emailEmptyErr    = "Email cannot be empty"

	cmdNotFoundErr = "Command not found in user module"

	notLoginErr = "Not Login"

	tokenByteSize = 36
	lenMin        = 7
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
		res.Error(cmdNotFoundErr)
	}
}

/*
 * ---
 */

func (svc *Service) reg(res *ws.Response, req *ws.Request) {
	email := req.ParamStr("email")
	if svc.repo.hasEmail(email) {
		res.Error(emailExistsErr + ": " + email)
		return
	}

	username := req.ParamStr("username")
	if username == "" {
		res.Error(usernameEmptyErr)
		return
	}
	if len(username) < 6 {
		res.Error(usernameTooShortErr)
		return
	}
	if svc.repo.hasUsername(username) {
		res.Error(usernameEmptyErr + ": " + username)
	}

	password := req.ParamStr("password")
	if len(password) < 7 {
		res.Error(passwordTooShortErr)
		return
	}

	passhash, err := hashPassword(password)
	if err != nil {
		res.Error(err.Error())
		return
	}

	uid, err := uuid.NewBase58()
	if err != nil {
		res.Error(err.Error())
		return
	}

	err = svc.repo.saveUser(uid, email, username, passhash)
	if err != nil {
		res.Error(err.Error())
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
		res.Error(emailEmptyErr)
		return
	}

	password := req.ParamStr("password")
	if password == "" {
		res.Error(passwordEmptyErr)
		return
	}

	uid := svc.repo.fetchUidByEmail(email)
	if uid == "" {
		res.Error(emailNotFoundErr + ": " + email)
		return
	}

	usr, err := svc.repo.fetchUser(uid)
	if err != nil {
		res.Error(err.Error())
		return
	}

	if !checkPasswordHash(password, usr.Passhash) {
		res.Error(passwordNotMatchErr)
		return
	}

	req.SetUid(uid)
}

func (svc *Service) logout(res *ws.Response, req *ws.Request) {
	req.RemoveUid()
}
