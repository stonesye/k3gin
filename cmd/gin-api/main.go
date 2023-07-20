package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"k3gin/app"
	"os"
)

func main() {

	// 将 tag 封装到 context 里面
	ctx := app.NewTagContext(context.Background(), "__main__")

	cliApp := cli.NewApp()
	cliApp.Name = "k3-gin"
	cliApp.Version = "1.0.1"
	cliApp.Usage = "K3-GIN based on gin + gorm + wire + logrus + rotatelogs."
	cliApp.Commands = []*cli.Command{
		webCmd(ctx),
	}

	// 命令行包cli 的Run函数 其实是执行 Commands下所有的 cli.Command 中的 Key =  Action 指定函数
	if err := cliApp.Run(os.Args); err != nil {
		app.WithContext(ctx).Errorf(err.Error())
	}
}

func webCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "api",
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
			return app.Run(ctx, app.SetConfigFile(c.String("conf")), app.SetWWWDir(c.String("www")))
		},
	}
}

/**
go run cmd/gin-admin/main.go web -c ./configs/config.toml --www ./static

swag init --parseDependency --generalInfo ./cmd/${APP}/main.go --output ./internal/app/swagger

# Or use Makefile: make swagger

wire gen ./internal/app

# Or use Makefile: make wire
*/
