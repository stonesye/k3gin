package cron

import (
	"context"
	"fmt"
	"github.com/google/wire"
	v3cron "github.com/robfig/cron/v3"
	"k3gin/app/cache/redisx"
	"k3gin/app/config"
	cronctx "k3gin/app/cron/context"
	"k3gin/app/cron/middleware"
	"k3gin/app/gormx"
	"k3gin/app/httpx"
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

type Cron struct {
	V3Cron      *v3cron.Cron
	Middlewares []cronctx.HandleFunc
	Db          *gormx.DB
	Redis       *redisx.Store
	HttpClient  *httpx.Client
}

var CronSet = wire.NewSet(wire.Struct(new(Cron), "Db", "Redis", "HttpClient"), gormx.InitGormDB, redisx.RedisStoreSet, httpx.InitHttp)

func (cron *Cron) waitGraceExit() int {
	stat := 0

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
			case <-cron.V3Cron.Stop().Done(): // 等待所有的定时任务结束
				stat = 0
			}
			return stat // 退出函数
		case syscall.SIGHUP:
			stat = 1
		default:
		}
	}

	return stat
}

func (cron *Cron) Use(handleFunc cronctx.HandleFunc) {
	cron.Middlewares = append(cron.Middlewares, handleFunc)
}

func (cron *Cron) registerMiddleware() {
	cron.Use(middleware.RecoveryCron())
}

func Run(ctx context.Context, opts ...Option) error {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// # 初始化config #
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	config.MustLoad(o.conf, o.cron)
	logger.WithContext(ctx).Printf("Start #CRON# server, #run_mode %s,#version %s,#pid %d", config.C.RunMode, o.version, os.Getpid())

	// # 初始化CRON #
	cron, cleanFunc, err := InitCron(ctx)
	if err != nil {
		return err
	}

	cron.V3Cron.Start()
	stat := cron.waitGraceExit()

	cleanFunc()

	logger.WithContext(ctx).Info("Cron Server exit !")
	os.Exit(stat)

	return nil
}

func InitCron(ctx context.Context) (*Cron, func(), error) {

	// 初始化 logrus
	loggerCleanFunc, err := logger.InitLogger()
	if err != nil {
		return nil, loggerCleanFunc, err
	}
	// 初始化 cron
	cron, cleanFunc, err := BuildCronInject()
	if err != nil {
		return nil, cleanFunc, err
	}
	cron.Middlewares = make([]cronctx.HandleFunc, 0)
	cron.V3Cron = v3cron.New()

	return cron, cleanFunc, err
}
