package main

import (
	"log"
	"strings"
	"net/http"
	"github.com/tecposter/tec-server-go/internal/ws"
	"github.com/tecposter/tec-server-go/internal/user"
)

const (
	userModule = "user"
)

func extractModule(cmd string) string {
	dotIndex := strings.Index(cmd, ".")
	if dotIndex <= 0 {
		return ""
	}
	return cmd[0:dotIndex]
}

func handleWs(res *ws.Response, req *ws.Request) {
	switch extractModule(req.Cmd) {
	case userModule:
		user.HandleWs(res, req)
	default:
		res.Error("Module Not Found")
	}
}

func main() {
	http.Handle("/ws", ws.Handler(handleWs))

	if err := http.ListenAndServe(":8765", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
