package ws_api

import (
	"github.com/google/wire"
	"github.com/gorilla/websocket"
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

var TestApiSet = wire.NewSet(wire.Struct(new(Test), "*"), wire.Bind(new(ApiInterface), new(*Test)))

func (t *Test) TestApi(ctx *ws_context.WSContext) {
	conn, err := ctx.Upgrader.Upgrade(ctx.GinCtx.Writer, ctx.GinCtx.Request, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			logger.WithFieldsFromWSContext(ctx).Errorf("websocket error: %v", err)
		}
		return
	}

	t.write(ctx, conn)
	t.read(ctx, conn)
}

func (t *Test) write(ctx *ws_context.WSContext, conn *websocket.Conn) {

	defer func() {
		conn.Close()
	}()

	go func() {

	}()
}

func (t *Test) read(ctx *ws_context.WSContext, conn *websocket.Conn) {
	defer conn.Close()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
