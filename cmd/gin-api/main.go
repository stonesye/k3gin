package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"k3gin/app"
)

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

	app.SetLevel(app.DEBUG)
	app.SetFormatter("")
	logrus.Error("akkkkkk-------->")
	logrus.WithFields(logrus.Fields{"name": "12"}).Error("kkkkkk-->>>>>")

	ctx := app.NewTraceIDContext(context.Background(), "123")

	fmt.Println(ctx.Value(app.TraceIDKey))

}
