package grpcx

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"k3gin/app/config"
	"k3gin/app/contextx"
	"k3gin/app/logger"
	"net"
	"os"
	"os/signal"
	"strconv"
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

// WaitGraceExit 优雅退出
func WaitGraceExit(ctx context.Context) int {
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

// InitGRPCServer 初始化RPC服务器
func InitGRPCServer(ctx context.Context, registers ...func(*grpc.Server)) func() {
	var serv *grpc.Server
	cfg := config.C.GRPC

	addr := fmt.Sprintf("%s:%s", cfg.Host, strconv.Itoa(cfg.Port))
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		logger.WithContext(ctx).Fatalf("failed to listen: %v", err)
	}

	go func() {
		var opts []grpc.ServerOption

		// 如果是TLS 写入TLS参数
		if cfg.CertFile != "" && cfg.KeyFile != "" {
			creds, err := credentials.NewServerTLSFromFile(cfg.CertFile, cfg.KeyFile)
			if err != nil {
				logger.WithContext(ctx).Fatalf("fialed to generate credentials : %v", err)
			}
			opts = []grpc.ServerOption{grpc.Creds(creds)}
		}

		serv = grpc.NewServer(opts...)

		// 注册grpc的proto对象
		for _, register := range registers {
			register(serv)
		}

		if err := serv.Serve(lis); err != nil {
			logger.WithContext(ctx).Fatalf("failed to server : %v", err)
		}
	}()

	return func() {
		logger.WithContext(contextx.NewTag(ctx, "__grpc__")).Infof("Stop grpc.server !")
		serv.Stop()
	}
}

func Run(ctx context.Context, opts ...func(*options)) error {
	var o options

	for _, opt := range opts {
		opt(&o)
	}
	// 初始化config
	config.MustLoad(o.ConfigFile)
	config.PrintWithJSON()
	logger.WithContext(ctx).Printf("Start #GRPC# server, #run_mode %s,#version %s,#pid %d", config.C.RunMode, o.Version, os.Getpid())

	// 初始化logger
	logCleanFunc, err := logger.InitLogger()
	if err != nil {
		return err
	}

	// 初始化主要的组件, db, redis, http
	/**
	db, cleanFunc, err := gormx.InitGormDB()
	store, cleanFunc, err := redisx.InitRedisStore()
	client, cleanFunc, err := httpx.InitHttp()
	*/

	// 过滤器

	// 初始化grpc服务端TCP协议
	grpcServerCleanFunc := InitGRPCServer(ctx)
	// 优雅退出
	stat := WaitGraceExit(ctx)

	// 清理多余的数据
	logCleanFunc()
	grpcServerCleanFunc()
	logger.WithContext(ctx).Info("GRPC server will been exit !")
	os.Exit(stat)
	return nil
}
