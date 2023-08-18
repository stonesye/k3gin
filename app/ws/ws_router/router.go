package ws_router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/gorilla/websocket"
	"io"
	"k3gin/app/logger"
	"k3gin/app/ws/ws_api"
	"k3gin/app/ws/ws_context"
	"net/http"
)

type IRouter interface {
	Register(*gin.Engine) error
}

type WSRouter struct {
	ws_api.Test
}

var WSRouterSet = wire.NewSet(wire.Struct(new(WSRouter), "*"), wire.Bind(new(IRouter), new(*WSRouter)))

func (w *WSRouter) Register(engine *gin.Engine) error {

	g := engine.Group("/ws")
	{
		v1 := g.Group("/v1")
		{
			v1.GET("", WithWSContext(w.Test.TestApi))
		}
	}

	return nil
}

func WithWSContext(handler func(*ws_context.WSContext)) func(*gin.Context) {
	return func(c *gin.Context) {
		// 每个链接都应该有独立的ctx
		ctx := ws_context.WSContext{
			Context: context.TODO(),
			GinCtx:  c,
			Upgrader: &websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
			KV: make(map[string]interface{}),
		}
		handler(&ctx)
	}
}

type IWSProxy interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Copy(io.Writer, io.Reader) (int64, error)
}

type WSProxy struct {
	ctx    *ws_context.WSContext
	ws     *websocket.Conn
	wsType int
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
	proxy.wsType, p, err = proxy.ws.ReadMessage()
	if err != nil {
		if err != io.EOF {
			logger.WithFieldsFromWSContext(proxy.ctx).Errorf("Read message error: %v", err)
		}
		return 0, err
	}
	return len(p), err
}

func (proxy *WSProxy) Write(p []byte) (n int, err error) {
	err = proxy.ws.WriteMessage(proxy.wsType, p)
	if err != nil {
		if err != io.EOF {
			logger.WithFieldsFromWSContext(proxy.ctx).Errorf("Writer message error: %v", err)
		}
		return 0, err
	}
	return len(p), err
}

func (proxy *WSProxy) Copy(writer io.Writer, reader io.Reader) (int64, error) {
	// 一次性读取4096  字节
	buf := make([]byte, 4096)

	var written int64 = 0

	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if n == 0 {
			return written, nil
		}
		written += int64(n)
		m, err := writer.Write(buf[:n])
		if err != nil {
			return written, err
		}

		if m > n {
			return written, io.ErrShortWrite
		}
	}
}
