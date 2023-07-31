package cron

import (
	"context"
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
	V3Cron      *v3cron.Cron
	Middlewares []cronctx.HandleFunc // 每个Task/Job都需要执行的全局中间件
	Db          *gormx.DB
	Redis       *redisx.Store
	HttpClient  *httpx.Client
}

var CronSet = wire.NewSet(wire.Struct(new(Cron), "Db", "Redis", "HttpClient"), gormx.InitGormDB, redisx.RedisStoreSet, httpx.InitHttp)

func (cron *Cron) waitGraceExit(ctx context.Context) int {
	stat := 0

	// 创建新号源， 控制cron的运行， 确保只有接触到特殊的信号以后， 主协程才会退出，子协程才会被回收
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	for {
		logger.WithContext(ctx).Info("Waiting signal exiting cron ... ")

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

func (cron *Cron) Use(handleFunc cronctx.HandleFunc) {
	cron.Middlewares = append(cron.Middlewares, handleFunc)
}

func (cron *Cron) registerMiddleware() {
	cron.Use(middleware.RecoveryCron())
}

// AddJob 每个Job都需要生成一个新的context ，每个context都只存储一个时间节点要执行的所有middleware, 每个时间节点可能要执行多个middleware
func (cron *Cron) AddJob(name string, spec string, handles ...cronctx.HandleFunc) error {
	middlewares := make([]cronctx.HandleFunc, len(cron.Middlewares)+len(handles))
	copy(middlewares, cron.Middlewares)
	copy(middlewares[len(cron.Middlewares):], handles)

	f := func() {
		ctx := cronctx.NewCronContext(name, spec, middlewares)
		ctx.Next()
	}

	job := v3cron.SkipIfStillRunning(v3cron.DefaultLogger)(newJob(f))
	_, err := cron.V3Cron.AddJob(spec, job)

	return err
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
	cron, cleanFunc, err := InitCron(ctx)
	if err != nil {
		return err
	}

	// # 注册中间件
	cron.registerMiddleware()

	// # 添加定时任务
	cron.AddJob("Task1", "*/2 * * * * ", task1)
	cron.AddJob("Task2", "*/2 * * * *", middleware.TimeoutCron(5*time.Second), task2)

	// # goroutine 执行定时任务
	cron.V3Cron.Start()

	// # 处理主协程的优雅的退出
	stat := cron.waitGraceExit(ctx)

	// # 清理垃圾信息
	loggerCleanFunc()
	cleanFunc()
	logger.WithContext(ctx).Info("Cron Server exited !")
	time.Sleep(time.Duration(1000) * time.Millisecond)
	os.Exit(stat)
	return nil
}

// InitCron 初始化 cron
func InitCron(ctx context.Context) (*Cron, func(), error) {
	cron, cleanFunc, err := BuildCronInject()

	if err != nil {
		return nil, cleanFunc, err
	}

	cron.Middlewares = make([]cronctx.HandleFunc, 0)
	cron.V3Cron = v3cron.New()
	return cron, cleanFunc, err
}
