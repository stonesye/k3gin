//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"k3gin/app/api"
	"k3gin/app/cache/redisx"
	"k3gin/app/dao/user"
	"k3gin/app/gormx"
	"k3gin/app/httpx"
	"k3gin/app/router"
	"k3gin/app/service"
)

func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		redisx.RedisStoreSet,
		httpx.InitHttp,
		gormx.InitGormDB,
		api.UserApiSet,
		service.UserSrvSet,
		user.UserRepoSet,
		router.RouterSet,
		router.InitGinEngine,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
