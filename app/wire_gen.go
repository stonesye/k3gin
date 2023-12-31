// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package app

import (
	"k3gin/app/api"
	"k3gin/app/cache/redisx"
	"k3gin/app/dao/user"
	"k3gin/app/gormx"
	"k3gin/app/grpcx"
	"k3gin/app/httpx"
	"k3gin/app/router"
	"k3gin/app/service"
)

// Injectors from wire.go:

func BuildInjector() (*Injector, func(), error) {
	db, cleanup, err := gormx.InitGormDB()
	if err != nil {
		return nil, nil, err
	}
	userRepo := &user.UserRepo{
		DB: db,
	}
	userSrv := &service.UserSrv{
		UserRepo: userRepo,
	}
	client, cleanup2, err := httpx.InitHttp()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	userApi := &api.UserApi{
		UserSrv: userSrv,
		Client:  client,
	}
	routerRouter := &router.Router{
		UserAPI: userApi,
	}
	engine := router.InitGinEngine(routerRouter)
	store, cleanup3, err := redisx.InitRedisStore()
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	clientConn, cleanup4, err := grpcx.InitClientRPC()
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	injector := &Injector{
		Engine:     engine,
		HttpClient: client,
		RedisStore: store,
		GrpcClient: clientConn,
	}
	return injector, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
