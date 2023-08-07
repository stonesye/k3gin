package grpcx

import (
	"context"
	"google.golang.org/grpc"
	"k3gin/app/config"
	"os"
	"os/signal"
	"syscall"
	"time"
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

type Server struct {
	gserver *grpc.Server
}

// WaitGraceExit 优雅退出
func WaitGraceExit(ctx context.Context, server *grpc.Server) int {
	var stat = 1
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	for {
		s := <-sig
		switch s {
		case syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM:
			select {
			//  等待一点时间关闭
			case <-time.NewTimer(time.Duration(config.C.GRPC.ShutdownTimeout) * time.Second).C:
				stat = 0
			}
			return stat
		case syscall.SIGHUP:
		default:
			return stat
		}
	}
}

func Run(ctx context.Context, options ...func(*options)) error {

	return nil
}
