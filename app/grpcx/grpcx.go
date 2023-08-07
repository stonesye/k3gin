package grpcx

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"k3gin/app/config"
	"k3gin/app/logger"
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

func Run(ctx context.Context, opts ...func(*options)) error {
	var server Server
	var o options

	for _, opt := range opts {
		opt(&o)
	}
	// 初始化config
	config.MustLoad(o.ConfigFile)
	// 初始化looger
	logger.InitLogger()

	// 初始化主要的组件

	// 初始化grpc服务端TCP协议

	// 优雅退出

	// 清理多余的数据

	stat := WaitGraceExit(ctx, server.gserver)
	fmt.Println(stat)

	return nil
}
