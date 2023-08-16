package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"k3gin/app/config"
	"k3gin/app/ws/ws_router"
)

var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"))

type Injector struct {
	Engine *gin.Engine
}

func initGinEngine(r ws_router.IRouter) *gin.Engine {
	gin.SetMode(config.C.RunMode)
	app := gin.New()
	r.Register(app)
	return app
}
