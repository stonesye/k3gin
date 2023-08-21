package ws_proxy

import (
	"github.com/gorilla/websocket"
	"io"
	"k3gin/app/logger"
	"k3gin/app/ws/ws_context"
)

type IWSProxy interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Copy(io.Writer, io.Reader) (int64, error)
}

type WSProxy struct {
	ctx *ws_context.WSContext
	ws  *websocket.Conn
}

func NewWebsocketConnect(ctx *ws_context.WSContext) (*WSProxy, func(), error) {
	ws, err := ctx.Upgrader.Upgrade(ctx.GinCtx.Writer, ctx.GinCtx.Request, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			logger.WithFieldsFromWSContext(ctx).Errorf("Fetch websocket conn error : %v", err)
		}
		return nil, func() {}, err
	}

	return &WSProxy{ws: ws, ctx: ctx}, func() {
		ws.Close()
	}, nil
}

func (proxy *WSProxy) Read(p []byte) (n int, err error) {
	_, nr, err := proxy.ws.NextReader()
	return nr.Read(p)
}

func (proxy *WSProxy) Write(p []byte) (n int, err error) {
	rw, _ := proxy.ws.NextWriter(websocket.TextMessage)
	defer rw.Close()
	return rw.Write(p)
}

func (proxy *WSProxy) Copy(writer io.Writer, reader io.Reader) (int64, error) {
	return io.Copy(writer, reader)
}
