package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"k3gin/app/config"
	"time"
)

func CORSMiddleware() gin.HandlerFunc {
	cfg := config.C.CORS

	// 利用 gin-contrib包来设置跨域
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           time.Second * time.Duration(cfg.MaxAge),
	})

}

/** 另外一种原生写法
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.GetHeader("Origin")
		c.Header("Access-Control-Allow-Origin", origin)                                                                                                                          // 注意这一行，不能配置为通配符“*”号
		c.Header("Access-Control-Allow-Credentials", "true")                                                                                                                     // 注意这一行，必须设定为 true
		c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Cookie, Origin, X-Requested-With, Content-Type, Accept, Authorization, Token, Timestamp, UserId") // 我们自定义的header字段都需要在这里声明
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type,cache-control")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			//c.AbortWithStatus(http.StatusNoContent)
			c.AbortWithStatus(http.StatusOK)
		}
		// 处理请求
		c.Next()
	}
}
*/
