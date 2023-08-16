package ws_router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/gorilla/websocket"
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
