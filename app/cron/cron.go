package cron

import (
	"context"
	"fmt"
	"github.com/google/wire"
	v3cron "github.com/robfig/cron/v3"
	"k3gin/app"
	"k3gin/app/config"
	"k3gin/app/logger"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
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

type Cronx struct {
	Cron *v3cron.Cron
}

var CronxSet = wire.NewSet(wire.Struct(new(Cronx), "*"))

func (cron *Cronx) waitGraceExit() {

	// 创建新号源， 控制cron的运行， 确保只有接触到特殊的信号以后， 主协程才会退出，子协程才会被回收
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	for {
		_, _ = fmt.Fprintln(os.Stdout, "等待CRON终止信号...")

		s := <-sig

		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			_, _ = fmt.Fprintf(os.Stdout, "收到信号: %s , CRON 服务正在退出...\n", s.String())
			select {
			case <-time.NewTimer(time.Duration(config.C.Cron.WaitGraceExit) * time.Millisecond).C: // 最多等待的时间
			case <-cron.Cron.Stop().Done(): // 等待所有的定时任务结束
			}
			return // 退出函数
		case syscall.SIGHUP:
		default:
		}
	}

}

func (cron *Cronx) Use() {

}

// Run 用于处理CronTab 的任务
func Run(ctx context.Context, opts ...Option) error {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cleanFunc, err := InitCron(ctx, opts...)
	if err != nil {
		return err
	}

	cleanFunc()
	return nil
}

func InitCron(ctx context.Context, opts ...Option) (func(), error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	// 初始化Config
	config.MustLoad(o.conf, o.cron)
	logger.WithContext(ctx).Printf("Start #CRON# server, #run_mode %s,#version %s,#pid %d", config.C.RunMode, o.version, os.Getpid())

	loggerCleanFunc, err := logger.InitLogger()
	if err != nil {
		return loggerCleanFunc, err
	}

	// 初始化
	cronx, cleanFunc, err := app.BuildCron()
	if err != nil {
		return cleanFunc, err
	}

	cronx.Cron.Run()
	cronx.waitGraceExit()

	return nil, nil
}
