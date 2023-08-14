package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"k3gin/app/config"
)

var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"))

type Injector struct {
	Engine *gin.Engine
}

func initGinEngine() *gin.Engine {
	gin.SetMode(config.C.RunMode)
	app := gin.New()
	return app
}
