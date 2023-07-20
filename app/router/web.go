package router

import (
	"github.com/gin-gonic/gin"
	"k3gin/app/config"
	"k3gin/app/middleware"
)

// InitGinEngine 生成一个Web 类型的 Gin Engine
func InitGinEngine(r IRouter) *gin.Engine {
	gin.SetMode(config.C.RunMode)

	app := gin.New()

	// 允许访问的目录地址
	prefixes := r.Prefixes()

	// 设置中间件
	if config.C.CORS.Enable {
		app.Use(middleware.CORSMiddleware())
	}

	if config.C.SESSION.Enable {
		app.Use(middleware.SESSMiddleware())
	}

	// 静态文件目录
	if dir := config.C.WWW; dir != "" {
		app.Use(middleware.WWWMiddleware(dir, middleware.AllowPathPrefixSkipper(prefixes...)))
	}

	// 将API封装进去, 如果有可能封装不同的api 可以改写r.Register
	if err := r.Register(app); err != nil {

	}

	return app
}
