//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	v3cron "github.com/robfig/cron/v3"
	"k3gin/app/api"
	"k3gin/app/cron"
	"k3gin/app/dao/user"
	"k3gin/app/gormx"
	"k3gin/app/router"
	"k3gin/app/service"
)

func BuildInjector() (*Injector, func(), error) {
	wire.Build(
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

func BuildCron(opts ...v3cron.Option) (*cron.Cronx, func(), error) {
	wire.Build(
		v3cron.New,
		cron.CronxSet,
		gormx.InitGormDB(),
	)
	return new(cron.Cronx), nil, nil
}
