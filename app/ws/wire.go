//go:build wireinject
// +build wireinject

package ws

func BuildInjector() (*Injector, func(), error) {
	return &Injector{}, nil, nil
}
