package router

import (
	"github.com/gin-gonic/gin"
	"k3gin/app/config"
	"k3gin/app/middleware"
)

// InitGinEngine 生成一个Web 类型的 Gin Engine
func InitGinEngine() *gin.Engine {
	gin.SetMode(config.C.RunMode)

	app := gin.New()

	// 设置中间件
	if config.C.CORS.Enable {
		app.Use(middleware.CORSMiddleware())
	}

	if config.C.SESSION.Enable {
		app.Use(middleware.SESSMiddleware())
	}

	return app
}
