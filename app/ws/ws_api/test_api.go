package ws_api

import (
	"github.com/google/wire"
	"github.com/gorilla/websocket"
	"io"
	"k3gin/app/cache/redisx"
	"k3gin/app/gormx"
	"k3gin/app/httpx"
	"k3gin/app/logger"
	"k3gin/app/ws/ws_context"
)

/**
测试websocket
*/

type Test struct {
	DB         *gormx.DB
	HttpClient *httpx.Client
	Redis      *redisx.Store
}

var TestApiSet = wire.NewSet(wire.Struct(new(Test), "*"), wire.Bind(new(WSApi), new(*Test)))

func (t *Test) TestApi(ctx *ws_context.WSContext) {
	conn, err := ctx.Upgrader.Upgrade(ctx.GinCtx.Writer, ctx.GinCtx.Request, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			logger.WithFieldsFromWSContext(ctx).Errorf("websocket error: %v", err)
		}
		return
	}

	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			logger.WithFieldsFromWSContext(ctx).Errorf("Failed to read message from ws : %v", err)
		}

		switch messageType {
		case websocket.TextMessage:
			err := conn.WriteMessage(messageType, p)
			break

		case websocket.BinaryMessage:
			err := conn.WriteMessage(messageType, p)
			break
		}
	}
}

func (t *Test) Read(p []byte) (n int, err error) {
	panic("implement me")
}

func (t *Test) Write(p []byte) (n int, err error) {
	panic("implement me")
}

func (t *Test) Copy(writer io.Writer, reader io.Reader) (int64, error) {
	panic("implement me")
}
