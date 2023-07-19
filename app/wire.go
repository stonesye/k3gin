//go:build wireinject
// +build wireinject

package app

import "github.com/google/wire"

/**
func TestInjector(name string, age int) (*util.Person, func(), error) {
	wire.Build(util.PSet)

	return new(util.Person), nil, nil
}
*/

func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		InitGormDB,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
