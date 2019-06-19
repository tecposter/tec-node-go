package main

import (
	"flag"
	"log"
	"path"
	"strings"
	"net/http"
	"github.com/tecposter/tec-server-go/internal/ws"
	"github.com/tecposter/tec-server-go/internal/iotool"
	"github.com/tecposter/tec-server-go/internal/user"
)

const (
	userModule = "user"

	dirMode = 0777
)

type application struct {
	userWsHdl *user.WsHandler
}

func (app *application) handleWs(res *ws.Response, req *ws.Request) {
	switch extractModule(req.Cmd) {
	case userModule:
		app.userWsHdl.Handle(res, req)
	default:
		res.Error("Module Not Found")
	}
}

func (app *application) Close() {
	app.userWsHdl.Close()
}


func main() {
	addr := ":8765"
	dataDir := getDataDir()

	log.Printf("data directory: %s, Serve addr: %s", dataDir, addr)


	app := &application{
		userWsHdl: getUserWsHdl(dataDir)}
	defer app.Close()

	http.Handle("/ws", ws.Handler(app.handleWs))
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getUserWsHdl(dataDir string) *user.WsHandler {
	hdl, err := user.NewWsHandler(path.Join(dataDir, "user"))
	if err != nil {
		log.Fatal(err)
	}
	return hdl
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
