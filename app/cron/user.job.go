package cron

import (
	"fmt"
	"time"
)

var count = 0

func UserJob(ctx *FrameContext) {
	count++
	fmt.Println("count = ", count)
	if count == 2 {
		panic("异常错误")
	}
	fmt.Println("this is user job", ctx.Ctx, ctx.Cron)
	time.Sleep(3 * time.Second)
}

func UserJobTimeout(ctx *Context) {
	<-ctx.Done()
	fmt.Println(ctx)
}
