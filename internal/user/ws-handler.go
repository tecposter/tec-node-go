package user

import (
	"log"
	//"bytes"
	//"encoding/binary"
	//"encoding/binary"
	"github.com/tecposter/tec-server-go/internal/ws"
	"github.com/tecposter/tec-server-go/internal/com/uuid"
	"github.com/tecposter/tec-server-go/internal/com/rand"
)

type WsHandler struct {
	repo *repository
	cache *cacheStorage
}

const (
	regCmd = "user.reg"
	loginCmd = "user.login"
	refreshTokenCmd = "user.refresh-token"

	usernameEmptyErr = "Usernaame cannot be empty"
	usernameTooShortErr = "Username too short - minimum length is 6"
	usernameExistsErr = "Username already exists"

	passwordTooShortErr = "Password too short - minimum length is 7"
	passwordNotMatchErr = "Password not match"
	passwordEmptyErr = "password cannot be empty"

	emailExistsErr = "Email already exists"
	emailFormatErr = "Error eamil format"
	emailNotFoundErr = "Email not found"
	emailEmptyErr = "Email cannot be empty"

	cmdNotFoundErr = "Command not found in user module"

	notLoginErr = "Not Login"

	tokenByteSize = 36
	lenMin = 7
)

func NewWsHandler(userDataDir string) (*WsHandler, error) {
	repo, err := newRepo(userDataDir)
	if err != nil {
		return nil, err
	}

	cache := newCacheStorage()
	wsHandler := &WsHandler{
		repo: repo,
		cache: cache}

	return wsHandler, nil
}

func (hdl *WsHandler) Close() {
	hdl.repo.Close()
}

func (hdl *WsHandler) Handle(res *ws.Response, req *ws.Request) {
	log.Printf("ws.Request: %+v\n", req)
	switch req.Cmd() {
	case regCmd:
		hdl.reg(res, req)
	case loginCmd:
		hdl.login(res, req)
	case refreshTokenCmd:
		hdl.refreshToken(res, req)
	default:
		res.Error(cmdNotFoundErr)
	}
}

func (hdl *WsHandler) reg(res *ws.Response, req *ws.Request) {
	email := req.ParamStr("email")
	if hdl.repo.hasEmail(email) {
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
	if hdl.repo.hasUsername(username) {
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

	err = hdl.repo.saveUser(uid, email, username, passhash)
	if err != nil {
		res.Error(err.Error())
		return
	}

	res.Load(map[string]interface{}{
		"uid": uid,
		"email": email,
		"username": username})
}

func (hdl *WsHandler) login(res *ws.Response, req *ws.Request) {
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

	uid := hdl.repo.fetchUidByEmail(email)
	if uid == "" {
		res.Error(emailNotFoundErr + ": " + email)
		return
	}

	usr, err := hdl.repo.fetchUser(uid)
	if err != nil {
		res.Error(err.Error())
		return
	}


	if !checkPasswordHash(password, usr.Passhash) {
		res.Error(passwordNotMatchErr)
		return
	}

	req.SetUid(uid)

	token, err := rand.GenerateStr(tokenByteSize)
	if err != nil {
		res.Error(err.Error())
		return
	}

	hdl.cache.set("token-" + uid, token)
	res.Set("token", token)
}

func (hdl *WsHandler) refreshToken(res *ws.Response, req *ws.Request) {
	uid := req.GetUid()
	if uid == "" {
		res.Error(notLoginErr)
	}

	token, err := rand.GenerateStr(tokenByteSize)
	if err != nil {
		res.Error(err.Error())
		return
	}
	hdl.cache.set("token-" + uid, token)
	res.Set("token", token)
}

/*
func (hdl *WsHandler) reg(res *ws.Response, req *ws.Request) {
	email := req.ParamStr("email")
	username := req.ParamStr("username")
	password := req.ParamStr("password")
	passhash, err := hashPassword(password)

	if err != nil {
		res.Error(err.Error())
		return
	}

	if hdl.st.Has([]byte(email)) {
		v, err := hdl.st.Get([]byte(email))
		if err != nil {
			res.Error(err.Error())
			return
		}

		checkByteVal(v)

		res.Error("email already exists: " + email)
		return
	}

	u := User {
		Email: email,
		Username: username,
		Passhash: passhash}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, &u)

	if err != nil {
		log.Println("binary.Write failed: ", err)
	}
	hdl.st.Set([]byte(email), buf.Bytes())

	log.Println("user ", u)
	log.Printf("%s, %s, %s, %s", email, username, password, passhash)
}

func checkByteVal(b []byte) {
	u := User{}

	r := bytes.NewReader(b)
	if err := binary.Read(r, binary.BigEndian, &u); err != nil {
		log.Println("binary.Read failed:", err)
	}

	log.Println(u)
}
*/
