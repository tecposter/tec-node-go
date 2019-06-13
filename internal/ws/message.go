package ws

import (
	"golang.org/x/net/websocket"
	"fmt"
)

type WsHandleFunc func(*Response, *Request)

func Handler(hdl WsHandleFunc) websocket.Handler {
	return websocket.Handler(func(conn *websocket.Conn) {
		for {
			var txt string
			if err := websocket.Message.Receive(conn, &txt); err != nil {
				fmt.Println("Can't receive")
				break
			}

			req, err := ParseRequest(txt)
			if err != nil {
				fmt.Println("ParseRequest: ", err)
			}

			res := NewResponse(req.Cmd)
			hdl(res, req)
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
