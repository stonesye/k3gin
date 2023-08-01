package cron

import (
	"context"
	"github.com/google/wire"
	v3cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"k3gin/app/cache/redisx"
	"k3gin/app/config"
	"k3gin/app/cron/job"
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
	version string // cron的version
}

func WithConf(conf string) func(*options) {
	return func(o *options) {
		o.conf = conf
	}
}

func WithVersion(version string) func(*options) {
	return func(o *options) {
		o.version = version
	}
}

type Option func(*options)

type Cron struct {
	V3Cron     *v3cron.Cron
	Db         *gormx.DB
	Redis      *redisx.Store
	HttpClient *httpx.Client
}

var CronSet = wire.NewSet(wire.Struct(new(Cron), "Db", "Redis", "HttpClient"), gormx.InitGormDB, redisx.RedisStoreSet, httpx.InitHttp)

// waitGraceExit 优雅退出
func (cron *Cron) waitGraceExit(ctx context.Context) int {
	stat := 0
	// 创建新号源， 控制cron的运行， 确保只有接触到特殊的信号以后， 主协程才会退出，子协程才会被回收
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	logger.WithContext(ctx).Info("Waiting signal exiting cron ... ")

	for {

		s := <-sig
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			logger.WithContext(ctx).Infof("Received signal : %s, cron server exiting ...", s.String())

			select {
			case <-time.NewTimer(time.Duration(config.C.Cron.WaitGraceExit) * time.Millisecond).C: // 最多等待的时间
			case <-cron.V3Cron.Stop().Done(): // 等待所有的定时任务结束
				stat = 0
			}
			return stat // 退出函数

		case syscall.SIGHUP:
		default:
			stat = 1
			return stat
		}
	}
}

func (cron *Cron) withV3Cron(ctx context.Context) {
	cron.V3Cron = v3cron.New(v3cron.WithSeconds(), v3cron.WithLogger(v3cron.VerbosePrintfLogger(logrus.StandardLogger())))
}

func Run(ctx context.Context, opts ...Option) error {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// # 初始化config #
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	config.MustLoad(o.conf)
	config.PrintWithJSON()
	logger.WithContext(ctx).Printf("Start #CRON# server, #run_mode %s,#version %s,#pid %d", config.C.RunMode, o.version, os.Getpid())

	// 初始化 logrus
	loggerCleanFunc, err := logger.InitLogger()
	if err != nil {
		return err
	}

	// # 初始化CRON #
	cron, cleanFunc, err := BuildCronInject()
	cron.withV3Cron(ctx)

	cron.V3Cron.AddJob("*/5 * * * * *", &job.UserJob{Name: "yelei"})

	cron.V3Cron.Start()
	// #监听#
	stat := cron.waitGraceExit(ctx)
	// # 清理垃圾信息
	loggerCleanFunc()
	cleanFunc()
	logger.WithContext(ctx).Info("Cron Server exited !")
	time.Sleep(time.Duration(1000) * time.Millisecond)
	os.Exit(stat)
	return nil
}
