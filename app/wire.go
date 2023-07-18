//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	. "k3gin/app/util"
)

func InitMisson(name string) (Mission, error) {
	wire.Build(NewMonster, NewPlayer, NewMission)

	return Mission{}, nil
}

func InitEndingA(name string) (EndingA, error) {
	wire.Build(EndingASet)
	return EndingA{}, nil
}

func InitEndingB(name string) (EndingB, error) {
	wire.Build(EndingBSet)
	return EndingB{}, nil
}

func BuildInjector() (*DB, func(), error) {
	wire.Build(InitGormDB)
	return &DB{}, nil, nil
}
