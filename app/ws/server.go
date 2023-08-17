package ws

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"k3gin/app/config"
	"k3gin/app/logger"
	"k3gin/app/ws/ws_middleware"
	"k3gin/app/ws/ws_router"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type options struct {
	configFile string
	version    string
}

func WithConfigFile(configFile string) func(*options) {
	return func(o *options) {
		o.configFile = configFile
	}
}

func WithVersion(version string) func(*options) {
	return func(o *options) {
		o.version = version
	}
}

// Run WebSocket Server
func Run(ctx context.Context, opts ...func(*options)) error {

	var o options

	for _, opt := range opts {
		opt(&o)
	}

	config.MustLoad(o.configFile)
	config.PrintWithJSON()

	logger.WithFieldsFromWSContext(ctx).Printf("Start #Websocket#, #run_mode %s, #pid %d, #version %s", config.C.RunMode, os.Getpid(), o.version)

	loggerCleanFunc, err := logger.InitLogger()
	if err != nil {
		return err
	}

	// 初始化要用的各类组件
	injector, injectorCleanFunc, err := BuildInjector()
	if err != nil {
		return err
	}

	// 初始化websocket
	wsCleanFunc := initWebsocket(ctx, injector.Engine)
	// 处理优雅退出
	stat := waitGraceExit(ctx)

	loggerCleanFunc()
	injectorCleanFunc()
	wsCleanFunc()

	logger.WithFieldsFromWSContext(ctx).Info("Websocket server exit !")
	time.Sleep(time.Second)
	os.Exit(stat)

	return nil
}

func initGinEngine(r ws_router.IRouter) *gin.Engine {
	gin.SetMode(config.C.RunMode)
	app := gin.New()

	// 集中管理异常
	app.Use(ws_router.WithWSContext(ws_middleware.RecoveryMiddleware()))

	// TODO 过滤静态目录

	r.Register(app)

	return app
}

func initWebsocket(ctx context.Context, handler http.Handler) (cleanFunc func()) {

	var C = config.C.WebSocket

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", C.Host, strconv.Itoa(C.Port)),
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		var err error
		if C.KeyFile != "" && C.CertFile != "" {
			srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = srv.ListenAndServeTLS(C.CertFile, C.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	cleanFunc = func() {
		srv.SetKeepAlivesEnabled(false)
		ctx, cancel := context.WithTimeout(ctx, time.Duration(C.ShutdownTimeout)*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.WithFieldsFromWSContext(ctx).Errorf("Websocket shutdown err : %v", err.Error())
		}
	}
	return
}

func waitGraceExit(ctx context.Context) (stat int) {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

EXIT:
	for {
		s := <-sig

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			stat = 0
			break EXIT

		case syscall.SIGHUP:
		default:
			stat = 1
			break EXIT
		}
	}

	logger.WithFieldsFromWSContext(ctx).Infof("Websocket will been exit #stat:%v# .", stat)
	return stat
}
