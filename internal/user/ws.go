package user

import (
	"log"
	"github.com/tecposter/tec-server-go/internal/ws"
)

func HandleWs(res *ws.Response, req *ws.Request) {
	log.Printf("ws.Request: %+v\n", req)
}
