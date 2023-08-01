// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package cron

import (
	"k3gin/app/cache/redisx"
	"k3gin/app/cron/job"
	"k3gin/app/gormx"
	"k3gin/app/httpx"
)

// Injectors from wire.go:

func BuildCronInject() (*Cron, func(), error) {
	userJobName := _wireUserJobNameValue
	userJobSpec := _wireUserJobSpecValue
	client, cleanup, err := httpx.InitHttp()
	if err != nil {
		return nil, nil, err
	}
	db, cleanup2, err := gormx.InitGormDB()
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
	userJob := &job.UserJob{
		Name:  userJobName,
		Spec:  userJobSpec,
		Http:  client,
		DB:    db,
		Store: store,
	}
	worker := &Worker{
		UserJob: userJob,
	}
	cron := InitV3Cron(worker)
	cronCron := &Cron{
		V3Cron: cron,
	}
	return cronCron, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

var (
	_wireUserJobNameValue = job.UserJobName("user")
	_wireUserJobSpecValue = job.UserJobSpec("*/2 * * * * *")
)
