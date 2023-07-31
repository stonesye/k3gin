package cron

import "context"

type options struct {
	conf string // 基础配置信息
	cron string // 定时任务的相关配置
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

type Option func(*options)

// Run 用于处理 CRONTab的任务
func Run(ctx context.Context, opts ...Option) error {

	return nil
}
