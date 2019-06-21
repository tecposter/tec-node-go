package post

import (
	"github.com/tecposter/tec-server-go/internal/ws"
)

type WsHandler struct {
}

func NewWsHandler(dataDir string) (*WsHandler, error) {
	hdl := &WsHandler{
	}

	return hdl, nil
}

func (hdl *WsHandler) Close() {
}

func (hdl *WsHandler) Handle(res *ws.Response, req *ws.Request) {
	res.Set("echo", "post.Handle")
}
