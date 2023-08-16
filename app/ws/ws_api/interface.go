package ws_api

import (
	"github.com/gorilla/websocket"
	"k3gin/app/ws/ws_context"
)

type ApiInterface interface {
	write(*ws_context.WSContext, *websocket.Conn)
	read(*ws_context.WSContext, *websocket.Conn)
}
