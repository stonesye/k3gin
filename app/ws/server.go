package ws

import (
	"context"
	"k3gin/app/config"
	"k3gin/app/logger"
	"os"
	"os/signal"
	"syscall"
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

func WihtVerion(version string) func(*options) {
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

	// 处理优雅退出
	waitGraceExit(ctx)

	logger.WithContext(ctx).Info("Websocket will been stop ...")
	cleanFunc()
	return nil
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
