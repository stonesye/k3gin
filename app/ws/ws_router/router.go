package ws_router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"k3gin/app/logger"
	"k3gin/app/ws/ws_api"
	"k3gin/app/ws/ws_context"
	"net/http"
)

type IRouter interface {
	Register(*gin.Engine) error
}

type WSRouter struct {
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (w *WSRouter) Register(engine *gin.Engine) error {

	g := engine.Group("/ws")
	{
		v1 := g.Group("/v1")
		{
			v1.GET("", WithWSContext(ws_api.TestApi))
		}
	}

	return nil
}

func WithWSContext(handler func(*ws_context.WSContext)) func(*gin.Context) {
	return func(c *gin.Context) {
		wscontext := &ws_context.WSContext{
			Context: context.TODO(),
			GinCtx:  c,
		}

		ws, err := upgrade.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			if _, ok := err.(websocket.HandshakeError); !ok {
				logger.WithContext(c).Errorf("Websocket err : %v", err)
			}

			return
		}

		go func() {
			handler(wscontext)
		}()

		for {
			_, _, err = ws.ReadMessage()
			if err != nil {
				break
			}
		}

	}
}
