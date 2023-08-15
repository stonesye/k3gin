package ws_context

import (
	"context"
	"github.com/gin-gonic/gin"
)

type WSContext struct {
	context.Context
	GinCtx *gin.Context
	KV     map[interface{}]interface{}
}

func (w *WSContext) Set(k interface{}, v interface{}) {
	w.KV[k] = v
}

func (w *WSContext) Get(k interface{}) interface{} {
	return w.KV[k]
}

func (w *WSContext) Value(k interface{}) {
	w.Context.Value(k)
}

func (w *WSContext) WithContext(ctx context.Context, k interface{}, v interface{}) {
	context.WithValue(ctx, k, v)
}
