package ws_context

import (
	"context"
	"github.com/gin-gonic/gin"
)

type WSContext struct {
	context.Context
	GinCtx *gin.Context
	KV     map[string]interface{}
}

type IWSContext interface {
	Set(interface{}, interface{})
	Get(interface{}) interface{}
}

type (
	traceCTX struct{}
	tagCTX   struct{}
	stackCTX struct{}
)

func (w *WSContext) Set(k interface{}, v interface{}) {
	if s, ok := k.(string); ok {
		w.KV[s] = v
	} else {
		context.WithValue(w.Context, k, v)
	}
}

func (w *WSContext) Get(k interface{}) interface{} {

	if s, ok := k.(string); ok {
		return w.KV[s]
	}
	return w.Context.Value(k)
}

func (w *WSContext) value(k interface{}) interface{} {
	if s, ok := k.(string); ok {
		return w.KV[s]
	}
	return w.Context.Value(k)
}

func NewTag(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, tagCTX{}, tag)
}

func FromTag(ctx context.Context) (string, bool) {
	v := ctx.Value(tagCTX{})
	if s, ok := v.(string); ok {
		return s, s != ""
	}
	return "", false
}

func NewStack(ctx context.Context, stack string) context.Context {
	return context.WithValue(ctx, stackCTX{}, stack)
}

func FromStack(ctx context.Context) (string, bool) {
	v := ctx.Value(stackCTX{})
	if s, ok := v.(string); ok {
		return s, s != ""
	}

	return "", false
}

func NewTrace(ctx context.Context, trace string) context.Context {
	return context.WithValue(ctx, traceCTX{}, trace)
}

func FromTrace(ctx context.Context) (string, bool) {
	v := ctx.Value(traceCTX{})
	if s, ok := v.(string); ok {
		return s, s != ""
	}
	return "", false
}
