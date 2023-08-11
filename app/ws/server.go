package ws

import (
	"context"
	"k3gin/app/config"
	"k3gin/app/logger"
)

type options struct {
	configFile string
}

func WithConfigFile(configFile string) func(*options) {
	return func(o *options) {
		o.configFile = configFile
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

	// 处理优雅退出

	cleanFunc()
	return nil
}
