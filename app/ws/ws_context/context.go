package ws_context

import (
	"github.com/gin-gonic/gin"
	"k3gin/app/cron/context"
)

type WSContext struct {
	context.Context
	GinContext *gin.Context
}
