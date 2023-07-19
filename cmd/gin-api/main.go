package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"k3gin/app"
	"os"
)

func main() {

	ctx := app.NewTagContext(context.Background(), "__main__")

	cliApp := cli.NewApp()
	cliApp.Name = "k3-gin"
	cliApp.Version = "1.0.1"
	cliApp.Usage = "K3-GIN based on gin + gorm + wire + logrus + rotatelogs."
	cliApp.Commands = []*cli.Command{
		webCmd(ctx),
	}

	if err := cliApp.Run(os.Args); err != nil { // CLI 返回的ERRO 其实是 cmd中的action属性的函数所在的error
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
				Aliases:  []string{"c"},
				Usage:    "App configuration file(.json, .yaml, .toml)",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "www",
				Aliases: []string{"w"},
				Usage:   "App static directory",
			},
		},
		Action: func(c *cli.Context) error {
			fmt.Println(c.String("conf"), ":", c.String("c"), ":", c.String("www")) // 获取参数的信息
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

/**
func main() {
	/**
	miss, err := app.InitMisson("aaaaa")
	if err != nil {
		fmt.Printf("aaaaaaa")
	}
	miss.Start()

	a, err := app.InitEndingA("123")
	b, err := app.InitEndingB("321")
	a.Appear()
	b.Appear()

*/

/**
fmt.Printf(os.Getwd())
config.MustLoad("configs/config.toml")
config.PrintWithJSON()

*/

/**
app.SetLevel(app.DEBUG)
app.SetFormatter("")
logrus.Error("akkkkkk-------->")
logrus.WithFields(logrus.Fields{"name": "12"}).Error("kkkkkk-->>>>>")

ctx := app.NewTraceIDContext(context.Background(), "123")

fmt.Println(ctx.Value(app.TraceIDKey))

*/
