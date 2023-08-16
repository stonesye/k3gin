//go:build wireinject
// +build wireinject

package ws

import (
	"github.com/google/wire"
	"k3gin/app/cache/redisx"
	"k3gin/app/gormx"
	"k3gin/app/httpx"
	"k3gin/app/ws/ws_api"
	"k3gin/app/ws/ws_router"
)

func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		redisx.RedisStoreSet,
		httpx.InitHttp,
		gormx.InitGormDB,
		ws_api.TestApiSet,
		ws_router.WSRouterSet,
		initGinEngine,
		InjectorSet)
	return new(Injector), nil, nil
}
