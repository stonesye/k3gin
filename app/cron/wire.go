//go:build wireinject
// +build wireinject

package cron

import (
	"github.com/google/wire"
)

func BuildCronInject() (*Cron, func(), error) {
	wire.Build(
		CronSet,
	)
	return new(Cron), nil, nil
}
