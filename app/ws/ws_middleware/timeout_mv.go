package ws_middleware

import (
	"context"
	"k3gin/app/ws/ws_context"
	"time"
)

func TimeoutMiddleware() func(*ws_context.WSContext) {
	return func(ctx *ws_context.WSContext) {
		timeoutContext, cancel := context.WithTimeout(ctx, time.Duration(10)*time.Second)
		defer cancel()
		ctx.Context = timeoutContext
		ctx.GinCtx.Next()
	}
}
