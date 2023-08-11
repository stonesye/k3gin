package ws

import (
	"context"
	"crypto/tls"
	"fmt"
	"k3gin/app/config"
	"k3gin/app/logger"
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

	cleanFunc, err := logger.InitLogger()
	if err != nil {
		return err
	}

	// 初始化websocket
	InitWebsocket(ctx, handler)

	// 处理优雅退出
	waitGraceExit(ctx)

	logger.WithContext(ctx).Info("Websocket will been stop ...")
	cleanFunc()
	return nil
}

func InitWebsocket(ctx context.Context, handler http.Handler) (cleanFunc func()) {

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
			logger.WithContext(ctx).Errorf("Websocket shutdown err : %v", err.Error())
		}
	}
	return
}

func waitGraceExit(ctx context.Context) (stat int) {

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		s := <-sig

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			stat = 0
			return

		case syscall.SIGHUP:
		default:
			stat = 1
			return
		}
	}
}

type WS struct {
}

func (w *WS) ServerHTTP(resp http.ResponseWriter, req *http.Request) {

}

func (w *WS) AddPath(path string) {

}
