package ws

import (
	"fmt"
	"golang.org/x/net/websocket"
)

type WsHandleFunc func(*Response, *Request)

func Handler(callback WsHandleFunc) websocket.Handler {
	return websocket.Handler(func(conn *websocket.Conn) {

		hdl := newConnHandler(conn)

		for {
			var txt string
			if err := websocket.Message.Receive(conn, &txt); err != nil {
				fmt.Println("Can't receive")
				break
			}

			req, err := NewRequest(hdl, txt)
			if err != nil {
				fmt.Println("ParseRequest: ", err)
			}

			res := NewResponse(req.Cmd())
			callback(res, req)
			fmt.Printf("recv text: %s\n", txt)
			//fmt.Printf("recv: %s, websocket.Message.Request: %+v\n", txt, req)

			bs, err := res.Marshal()
			if err != nil {
				fmt.Println("Response.Marshal: ", err)
				break
			}

			if err = websocket.Message.Send(conn, string(bs)); err != nil {
				fmt.Println("websocket.Message.Send: ", err)
				break
			}
		}
	})
}
