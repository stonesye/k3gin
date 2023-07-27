package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"k3gin/app/config"
	"k3gin/app/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// options 封装命令行输入的参数值
type options struct {
	ConfigFile string // 配置文件地址
	WWWDir     string // 静态文件地址
}

func SetConfigFile(s string) func(*options) {
	return func(o *options) {
		o.ConfigFile = s
	}
}

func SetWWWDir(s string) func(*options) {
	return func(o *options) {
		o.WWWDir = s
	}
}

func Run(ctx context.Context, opts ...func(*options)) error {
	// 初始化程序退出状态
	state := 1
	// 创建一个信号chan
	sc := make(chan os.Signal, 1)
	// 设置允许传递给 singal chan 的信号类型
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 初始化所有的对象
	cleanFunc, err := Init(ctx, opts...)

	if err != nil {
		return err
	}

EXIT:
	select {
	case sig := <-sc:
		// 打印 signal chan接收到的信号
		logger.WithContext(ctx).Infof("Receive signal[%s]", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			// 如果接收到的信号是以上信号，忽略钓, 重新接收信号
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			// 如果信号是sighup, 或者其他的 则退出
			state = 1
			break EXIT
		}
	}

	// 清理程序，并退出
	cleanFunc()
	logger.WithContext(ctx).Info("Server exit !")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}

// Init  (ctx [自定义context 里面封装了一些Key-value 区别于gin.Context] , opts [用于初始化options Struct])
func Init(ctx context.Context, opts ...func(*options)) (func(), error) {
	var o options

	// 初始化CLI传递的配置文件信息，封装到options struct
	for _, opt := range opts {
		opt(&o)
	}

	// 加载config文件内容到Config strut
	config.MustLoad(o.ConfigFile)

	// 启动命令会可选的方式带入静态目录， 如果没有附带就沿用配置文件的值
	if v := o.WWWDir; v != "" {
		config.C.WWW = v
	}

	// 检查一下所有的配置文件, 打印
	config.PrintWithJSON()

	// 利用默认的logrus来打印日志, 并没有利用到定制化的logrus, 因为还没有调用InitLogger
	logger.WithContext(ctx).Printf("Start server,#run_mode %s,#pid %d", config.C.RunMode, os.Getpid())

	// 初始化logrus 定制化日志
	loggerCleanFunc, err := logger.InitLogger()
	if err != nil {
		return nil, err
	}

	// 利用wire初始化所有的类
	injector, injectorCleanFunc, err := BuildInjector()
	if err != nil {
		return nil, err
	}

	// 开协程 监听HTTP/HTTPS
	httpServerCleanFunc := InitHttpServer(ctx, injector.Engine)

	return func() {
		httpServerCleanFunc()
		loggerCleanFunc()
		injectorCleanFunc()
	}, nil
}

// InitHttpServer 初始化HTTP服务器
func InitHttpServer(ctx context.Context, handler http.Handler) func() {
	cfg := config.C.HTTP

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		logger.WithContext(ctx).Printf("HTTP server is running at %s.", addr)

		var err error

		if cfg.CertFile != "" && cfg.KeyFile != "" {
			srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// 返回一个函数，在这个函数里面可以做httpserver的清理，方便整个应用退出后，清理工作进行
	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.ShutdownTimeout))
		defer cancel()

		// 将长连接关闭掉
		srv.SetKeepAlivesEnabled(false)
		// 关闭HTTPServer服务, 设置了timeout context  一旦时间到了就会执行ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			logger.WithContext(ctx).Errorf(err.Error())
		}
	}
}
