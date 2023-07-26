package middleware

import (
	"github.com/gin-gonic/gin"
	"k3gin/app/contextx"
	"k3gin/app/util/trace"
)

func TraceMiddleware(prefixes func(*gin.Context) bool) gin.HandlerFunc {

	return func(c *gin.Context) {

		if prefixes(c) {
			c.Next()
			return
		}

		traceID := c.GetHeader("X-Request-Id")

		if traceID == "" {
			traceID = trace.NewTraceID()
		}

		ctx := contextx.NewTraceID(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)

		c.Writer.Header().Set("X-Trace-Id", traceID)
		c.Next()
	}
}
