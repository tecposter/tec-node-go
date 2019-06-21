package main

import (
	"flag"
	"log"
	"path"
	"strings"
	"net/http"
	"github.com/tecposter/tec-server-go/internal/com/iotool"

	"github.com/tecposter/tec-server-go/internal/ws"

	"github.com/tecposter/tec-server-go/internal/user"
	"github.com/tecposter/tec-server-go/internal/post"
)

const (
	userModule = "user"
	postModule = "post"

	bindAddrDefault = ":8765"

	dirMode = 0777

	notLoginErr = "Not Login"
	moduleNotFoundErr = "Module not found"
)

/*
 * main
 */

func main() {
	dataDir := getDataDir()
	bindAddr := getBindAddr()
	log.Printf("data directory: %s, binding addr: %s", dataDir, bindAddr)

	app := getApp(dataDir)
	defer app.Close()

	http.Handle("/ws", ws.Handler(app.handleMsg))
	if err := http.ListenAndServe(bindAddr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

/*
 * application
 */

type application struct {
	postWsHdl *post.WsHandler
	userWsHdl *user.WsHandler
}

func (app *application) handleMsg(res *ws.Response, req *ws.Request) {
	mdl := extractModule(req.Cmd())

	switch mdl {
	case userModule:
		app.userWsHdl.Handle(res, req)
	case postModule:
		requireLogin(res, req, app.postWsHdl.Handle)
	default:
		res.Error(moduleNotFoundErr + ": " + mdl)
	}
}

func (app *application) Close() {
	app.userWsHdl.Close()
	app.postWsHdl.Close()
}

func getApp(dataDir string) *application {
	return &application{
		postWsHdl: getPostWsHdl(dataDir),
		userWsHdl: getUserWsHdl(dataDir)}
}

func getUserWsHdl(dataDir string) *user.WsHandler {
	hdl, err := user.NewWsHandler(path.Join(dataDir, "user"))
	if err != nil {
		log.Fatal(err)
	}
	return hdl
}

func getPostWsHdl(dataDir string) *post.WsHandler {
	hdl, err := post.NewWsHandler(dataDir)
	if err != nil {
		log.Fatal(err)
	}
	return hdl
}

/*
 * local func
 */

func requireLogin(res *ws.Response, req *ws.Request, callback ws.HandleMsgFunc) {
	if req.GetUid() == "" {
		res.Error(notLoginErr)
		return
	}
	callback(res, req)
}

func getBindAddr() string {
	bindAddr := flag.String("bind", bindAddrDefault, "Bind Addr")
	return *bindAddr
}

func getDataDir() string {
	currDir, err := iotool.CurrDir()
	if err != nil {
		log.Fatal(err)
	}
	dataDir := flag.String("datadir", path.Join(currDir, "data"), "Data Directory")
	iotool.MkdirIfNotExist(*dataDir)
	return *dataDir
}

func extractModule(cmd string) string {
	dotIndex := strings.Index(cmd, ".")
	if dotIndex <= 0 {
		return ""
	}
	return cmd[0:dotIndex]
}
