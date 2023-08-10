package cron

import (
	"context"
	"github.com/google/wire"
	v3cron "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"k3gin/app/cache/redisx"
	"k3gin/app/config"
	croncontext "k3gin/app/cron/context"
	"k3gin/app/cron/job"
	"k3gin/app/gormx"
	"k3gin/app/httpx"
	"k3gin/app/logger"
	"os"
	"os/signal"
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
	V3Cron         *v3cron.Cron
	HttpClient     *httpx.Client
	DB             *gormx.DB
	Store          *redisx.Store
	GrpcClient     *grpc.ClientConn
	GlobalJobFuncs []func(*croncontext.Context) // 所有Cron都需要执行的任务
}

var CronSet = wire.NewSet(wire.Struct(new(Cron), "V3Cron", "HttpClient", "DB", "Store", "GrpcClient"))

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

func (cron *Cron) AddJob(name string, spec string, jobs ...func(*croncontext.Context)) error {
	// # 将全局的和需要添加的任务都集中 #
	funcs := make([]func(*croncontext.Context), len(cron.GlobalJobFuncs)+len(jobs))
	copy(funcs, cron.GlobalJobFuncs)
	copy(funcs[len(cron.GlobalJobFuncs):], jobs)

	// 将任务封装成v3cron认识的job

	f := func() {
		// 当前任务 name， spec， funcs 元素都齐了, 创建job context
		ctx := croncontext.NewJobContext(name, spec, funcs...)
		ctx.Next() // 执行任务
	}

	_, err := cron.V3Cron.AddJob(spec, v3cron.SkipIfStillRunning(v3cron.DefaultLogger)(job.NewJob(f)))

	return err
}

func (cron *Cron) Use(jobFunc func(*croncontext.Context)) {
	if cron.GlobalJobFuncs == nil {
		cron.GlobalJobFuncs = make([]func(*croncontext.Context), 0)
	}
	cron.GlobalJobFuncs = append(cron.GlobalJobFuncs, jobFunc)
}

func (cron *Cron) WithFrameContext(handle func(*croncontext.FrameContext)) func(*croncontext.Context) {
	return func(ctx *croncontext.Context) {
		frameCtx := &croncontext.FrameContext{
			Context:     context.TODO(),
			HttpClient:  cron.HttpClient,
			DB:          cron.DB,
			Store:       cron.Store,
			GrpcClient:  cron.GrpcClient,
			CronContext: ctx,
		}
		handle(frameCtx)
	}
}

func InitV3Cron() *v3cron.Cron {
	return v3cron.New(v3cron.WithSeconds(), v3cron.WithLogger(v3cron.VerbosePrintfLogger(logrus.StandardLogger())))
}

func Run(ctx context.Context, opts ...Option) error {

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

	// # 设置所有Cron都需要执行的任务，比如异常处理 #
	cron.Use(job.RecoverGlobalJob())

	// # Add Job#
	_err := Register(cron)
	if _err != nil {
		return _err
	}

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

// Register 注册所有的Cron任务
func Register(cron *Cron) (_err error) {
	// cron.AddJob("userjob", "* * * * * *", cron.WithFrameContext(UserJob))
	// _err = cron.AddJob("userjob", "* * * * * *", cron.WithFrameContext(job.TestJob))
	// _err = cron.AddJob("Task2", "* * * * * *", job.TimeoutGlobalJob(5*time.Second), job.TestTimeoutJob)
	_err = cron.AddJob("GRPC", "*/5 * * * * *", cron.WithFrameContext(job.TestGRPC))
	return
}
