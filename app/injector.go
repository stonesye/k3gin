package app

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"k3gin/app/cache/redisx"
	"k3gin/app/httpx"
)

var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"), httpx.InitHttp, redisx.RedisStoreSet)

type Injector struct {
	Engine     *gin.Engine
	HttpClient *httpx.Client
	RedisStore redisx.Storer
}
