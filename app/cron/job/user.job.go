package job

import (
	"fmt"
	"k3gin/app/cron/context"
	"time"
)

var count = 0

func TestJob(ctx *context.FrameContext) {
	count++
	fmt.Println("count = ", count)
	if count == 2 {
		panic("异常错误")
	}
	fmt.Println("this is user job", ctx, ctx.CronContext)
	time.Sleep(3 * time.Second)
}

func TestTimeoutJob(ctx *context.Context) {
	<-ctx.Done()
	fmt.Println(ctx)
}
