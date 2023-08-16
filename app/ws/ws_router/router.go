package ws_router

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"k3gin/app/ws/ws_api"
	"k3gin/app/ws/ws_context"
)

type IRouter interface {
	Register(*gin.Engine) error
}

type WSRouter struct {
	upgrader *websocket.Upgrader
	ws_api.Test
}

/**
var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
*/

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
