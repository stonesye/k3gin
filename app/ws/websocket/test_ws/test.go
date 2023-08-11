package test_ws

import (
	"github.com/gorilla/websocket"
	"k3gin/app/ws"
	"log"
)

type test struct {
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (t *test) TestWS(ctx ws.Context) {
	ws, err := upgrader.Upgrade(ctx.Response, ctx.Request, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	go func() {
		write(ws)
	}()

	reader(ws)

}

func write(ws *websocket.Conn) {
	defer func() {
		ws.Close()
	}()
	if err := ws.WriteMessage(websocket.TextMessage, []byte("aaaaaaa")); err != nil {
		return
	}
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	_, _, err := ws.ReadMessage()
	if err != nil {
		return
	}
}
