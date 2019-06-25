package ws

import (
	"fmt"
	"golang.org/x/net/websocket"
)

type HandleMsgFunc func(*Response, *Request)
type HandleConnFunc func(*Connection)
type Handler interface {
	HandleConn(*Connection)
	HandleMsg(*Response, *Request)
}

func Handle(hdl Handler) websocket.Handler {
	return websocket.Handler(func(innerConn *websocket.Conn) {

		conn := newCollection(innerConn)
		hdl.HandleConn(conn)

		for {
			var txt string
			if err := websocket.Message.Receive(innerConn, &txt); err != nil {
				fmt.Println("Can't receive")
				break
			}

			req, err := NewRequest(conn, txt)
			if err != nil {
				fmt.Println("ParseRequest: ", err)
			}

			res := NewResponse(req.Cmd())
			hdl.HandleMsg(res, req)
			//callback(res, req)
			fmt.Printf("recv text: %s\n", txt)
			//fmt.Printf("recv: %s, websocket.Message.Request: %+v\n", txt, req)

			bs, err := res.Marshal()
			if err != nil {
				fmt.Println("Response.Marshal: ", err)
				break
			}

			if err = websocket.Message.Send(innerConn, string(bs)); err != nil {
				fmt.Println("websocket.Message.Send: ", err)
				break
			}
		}
	})
}
