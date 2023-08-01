//go:build wireinject
// +build wireinject

package cron

import (
	"github.com/google/wire"
	"k3gin/app/cache/redisx"
	"k3gin/app/gormx"
	"k3gin/app/httpx"
)

func BuildCronInject() (*Cron, func(), error) {
	wire.Build(
		gormx.InitGormDB,
		redisx.RedisStoreSet,
		httpx.InitHttp,
		InitV3Cron,
		CronSet,
	)
	return new(Cron), nil, nil
}
