package ws_router

import (
	"github.com/gin-gonic/gin"
	"k3gin/app/ws/ws_context"
)

type IRouter interface {
	Register(*gin.Engine) error
}

type WSRouter struct {
}

func (w *WSRouter) Register(engine *gin.Engine) error {

	g := engine.Group("/ws")
	{
		v1 := g.Group("/v1")
		{
			v1.GET("", func(ctx *gin.Context) {

			})
		}
	}

	return nil
}

func WithWSContext(handler func(ctx *ws_context.WSContext)) func(*gin.Context) {
	return func(ctx *gin.Context) {
		wscontext := &ws_context.WSContext{
			GinContext: ctx,
		}
		handler(wscontext)
	}
}
