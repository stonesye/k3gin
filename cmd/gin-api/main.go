package main

import (
	"context"
	_ "github.com/apache/skywalking-go"
	"github.com/urfave/cli/v2"
	"k3gin/app"
	"k3gin/app/contextx"
	"k3gin/app/cron"
	"k3gin/app/grpcx"
	"k3gin/app/logger"
	"k3gin/app/ws"
	_ "k3gin/cmd/gin-api/docs"
	"os"
	"runtime"
)

const VERSION = "1.0.1"

//	@title			k3gin
//	@version		1.0.1
//	@description	RBAC scaffolding based on GIN + GORM + WIRE.

//	@contact.name	STones_
//	@contact.email	yelei@3k.com
//	@contact.url	http://www.swagger.io/support

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @schemes	http https
// @host		127.0.0.1:8081
// @BasePath	/
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// 将 tag 封装到 自定义的 context 里面
	ctx := contextx.NewTag(context.Background(), "__main__")

	// 初始化CLI命令行对象
	cliApp := cli.NewApp()
	cliApp.Name = "k3-gin"
	cliApp.Version = VERSION
	cliApp.Usage = "K3-GIN based on gin + gorm + wire + logrus + rotatelogs + robfig-cron + swagger."
	cliApp.Commands = []*cli.Command{
		cmdWEB(ctx),
		cmdCRON(ctx),
		cmdRPC(ctx),
		cmdWS(ctx),
	}

	// 命令行包cli 的Run函数 其实是执行 Commands 下所有的 cli.Command 中的 Action 指定的函数
	if err := cliApp.Run(os.Args); err != nil {
		logger.WithContext(ctx).Errorf(err.Error())
	}
}

func cmdWS(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "ws",
		Usage: "Run websocket server group",
		Action: func(c *cli.Context) error {
			return ws.Run(ctx, ws.WithConfigFile(c.String("conf")), ws.WithVersion(VERSION))
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "App configuration file(.json, .yaml, .toml)",
				Required: true,
			},
		},
	}

}

func cmdRPC(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "rpc",
		Usage: "Run grpc server group",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "App configuration file(.json, .yaml, .toml)",
				Required: true,
			},
		},

		Action: func(c *cli.Context) error {
			return grpcx.Run(ctx, grpcx.WithVersion(VERSION), grpcx.WithConfigFile(c.String("conf")))
		},
	}
}

func cmdCRON(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "cron",
		Usage: "Run cron server group",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"}, // 别名
				Usage:    "App configuration file(.json, .yaml, .toml)",
				Required: true, // 是否一定要指定 --conf 或 -c
			},
		},
		Action: func(c *cli.Context) error {
			return cron.Run(ctx,
				cron.WithConf(c.String("conf")),
				cron.WithVersion(VERSION),
			)
		},
	}
}

func cmdWEB(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "Run http server group",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"}, // 别名
				Usage:    "App configuration file(.json, .yaml, .toml)",
				Required: true, // 是否一定要指定 --conf 或 -c
			},
			&cli.StringFlag{
				Name:    "www",
				Aliases: []string{"w"}, // 别名
				Usage:   "App static directory",
			},
		},
		Action: func(c *cli.Context) error {
			// c.String 获取命令行cli.StringFlag中指定的Name对应的Value
			return app.Run(ctx,
				app.SetConfigFile(c.String("conf")),
				app.SetWWWDir(c.String("www")),
				app.SetVersion(VERSION),
			)
		},
	}
}
