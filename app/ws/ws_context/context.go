package ws_context

import (
	"context"
	"github.com/gin-gonic/gin"
)

type WSContext struct {
	context.Context
	GinCtx *gin.Context
}
