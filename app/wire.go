//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"k3gin/app/api"
	"k3gin/app/dao/user"
	"k3gin/app/gormx"
	"k3gin/app/httpx"
	"k3gin/app/router"
	"k3gin/app/service"
)

/**
func TestInjector(name string, age int) (*util.Person, func(), error) {
	wire.Build(util.PSet)

	return new(util.Person), nil, nil
}
*/

func BuildInjector() (*Injector, func(), error) {
	wire.Build(
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
