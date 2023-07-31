package cron

import (
	"context"
	"k3gin/app/config"
	"k3gin/app/logger"
	"os"
)

type options struct {
	conf    string // 基础配置信息
	cron    string // 定时任务的相关配置
	version string // cron的version
}

func WithConf(conf string) func(*options) {
	return func(o *options) {
		o.conf = conf
	}
}

func WithCron(cron string) func(*options) {
	return func(o *options) {
		o.cron = cron
	}
}

func WithVersion(version string) func(*options) {
	return func(o *options) {
		o.version = version
	}
}

type Option func(*options)

// Run 用于处理CronTab 的任务
func Run(ctx context.Context, opts ...Option) error {
	var o options

	for _, opt := range opts {
		opt(&o)
	}

	// 初始化Config
	config.MustLoad(o.conf, o.cron)

	logger.WithContext(ctx).Printf("Start #CRON# server, #run_mode %s,#version %s,#pid %d", config.C.RunMode, o.version, os.Getpid())

	return nil
}
