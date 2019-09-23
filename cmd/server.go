package main

import (
	"errors"
	"flag"
	"github.com/tecposter/tec-node-go/internal/com/iotool"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/tecposter/tec-node-go/internal/ws"

	"github.com/tecposter/tec-node-go/internal/draft"
	"github.com/tecposter/tec-node-go/internal/post"
)

const (
	postModule  = "post"
	draftModule = "draft"
)

const (
	bindAddrDefault = "127.0.0.1:8765"
	dirMode         = 0777
)

// errors
var (
	ErrModuleNotFound = errors.New("Module not found")
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

	http.Handle("/ws", ws.Handle(app))
	if err := http.ListenAndServe(bindAddr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

/*
 * application
 */

type application struct {
	dataDir  string
	postSvc  *post.Service
	draftSvc *draft.Service
}

func getApp(dataDir string) *application {
	return &application{
		dataDir:  dataDir,
		postSvc:  getPostSvc(dataDir),
		draftSvc: getDraftSvc(dataDir)}
}

func (app *application) Close() {
	app.postSvc.Close()
	app.draftSvc.Close()
}

func (app *application) HandleMsg(res *ws.Response, req *ws.Request) {
	mdl := extractModule(req.Cmd())

	switch mdl {
	case draftModule:
		app.draftsvc.HandleMsg(res, req)
	case postModule:
		app.postSvc.HandleMsg(res, req)
	default:
		res.Error(ErrModuleNotFound)
	}
}

func (app *application) HandleConn(conn *ws.Connection) {
	conn.Set("dataDir", app.dataDir)
}

func getDraftSvc(dataDir string) *draft.Service {
	hdl, err := draft.NewService(dataDir)
	if err != nil {
		log.Fatal(err)
	}
	return hdl
}

func getPostSvc(dataDir string) *post.Service {
	hdl, err := post.NewService(dataDir)
	if err != nil {
		log.Fatal(err)
	}
	return hdl
}

/*
 * local func
 */

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
	flag.Parse()

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
