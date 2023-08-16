package ws_router

import (
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
	Upgrader *websocket.Upgrader
	ws_api.Test
}

func InitUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

var WSRouterSet = wire.NewSet(wire.Struct(new(WSRouter), "*"), wire.Bind(new(IRouter), new(*WSRouter)), InitUpgrader)

func (w *WSRouter) Register(engine *gin.Engine) error {

	g := engine.Group("/ws")
	{
		v1 := g.Group("/v1")
		{
			v1.GET("", w.WithWSContext(w.Test.TestApi))
		}
	}

	return nil
}

func (w *WSRouter) WithWSContext(handler func(*ws_context.WSContext)) func(*gin.Context) {
	return func(c *gin.Context) {

	}
}
