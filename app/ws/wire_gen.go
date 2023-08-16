// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package ws

import (
	"k3gin/app/cache/redisx"
	"k3gin/app/gormx"
	"k3gin/app/httpx"
	"k3gin/app/ws/ws_api"
	"k3gin/app/ws/ws_router"
)

// Injectors from wire.go:

func BuildInjector() (*Injector, func(), error) {
	upgrader := ws_router.InitUpgrader()
	db, cleanup, err := gormx.InitGormDB()
	if err != nil {
		return nil, nil, err
	}
	client, cleanup2, err := httpx.InitHttp()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	store, cleanup3, err := redisx.InitRedisStore()
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	test := ws_api.Test{
		DB:         db,
		HttpClient: client,
		Redis:      store,
	}
	wsRouter := &ws_router.WSRouter{
		Upgrader: upgrader,
		Test:     test,
	}
	engine := initGinEngine(wsRouter)
	injector := &Injector{
		Engine: engine,
	}
	return injector, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
