package grpcx

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type options struct {
	ConfigFile string // 配置文件地址
	Version    string
}

func WithConfigFile(conf string) func(*options) {
	return func(o *options) {
		o.ConfigFile = conf
	}
}

func WithVersion(version string) func(*options) {
	return func(o *options) {
		o.Version = version
	}
}

func WaitGraceExit(ctx context.Context) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

}

func Run(ctx context.Context, options ...func(*options)) error {
	return nil
}
