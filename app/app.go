package app

import (
	"context"
	"k3gin/app/config"
	"os"
)

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

func Init(ctx context.Context, opts ...func(*options)) (func(), error) {
	var o options

	// 初始化CLI传递的配置文件信息，封装到options struct
	for _, opt := range opts {
		opt(&o)
	}

	// 加载config文件内容到Config strut
	config.MustLoad(o.ConfigFile) // 并没有利用到定制化的logrus

	if v := o.WWWDir; v != "" {
		config.C.WWW = v
	}

	config.PrintWithJSON()

	// 利用默认的logrus来打印日志
	WithContext(ctx).Printf("Start server,#run_mode %s,#pid %d", config.C.RunMode, os.Getpid())

	// 初始化logrus 定制化日志
	loggerCleanFunc, err := InitLogger()
	if err != nil {
		return nil, err
	}

	// 利用wire初始化所有的类

	return func() {
		loggerCleanFunc()
	}, nil
}

func Run(ctx context.Context, opts ...func(*options)) error {

	_, err := Init(ctx, opts...)
	if err != nil {
		return err
	}

	return nil
}
