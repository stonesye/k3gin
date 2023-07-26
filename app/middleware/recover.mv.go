package middleware

import "github.com/gin-gonic/gin"

// RecoveryMiddleware 收集异常, 集中返回给Request端
func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

			}

		}()
	}
}
