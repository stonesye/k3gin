package job

import (
	"fmt"
	"k3gin/app/cron/context"
	"k3gin/app/grpcx/proto/test"
	"k3gin/app/logger"
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

func TestGRPC(ctx *context.FrameContext) {

	c := ctx.GrpcClient

	client := test.NewTestInfoClient(c)

	err := test.CallServerGetTestID(ctx, client, &test.Test{
		ID:      "1",
		Content: "Test RPC",
		Flag:    true,
	})

	if err != nil {
		logger.WithFieldsFromContext(ctx).Errorf("received error : %v", err)
	}

	err = test.CallServerStreamEcho(ctx, client, &test.TestRequest{Message: "RPC stream Testing "})
	if err != nil {
		logger.WithFieldsFromContext(ctx).Errorf("received error : %v", err)
	}

}
