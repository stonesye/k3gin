package ws_api

import (
	"fmt"
	"github.com/google/wire"
	"k3gin/app/cache/redisx"
	"k3gin/app/gormx"
	"k3gin/app/httpx"
	"k3gin/app/ws/ws_context"
)

/**
测试websocket
*/

type Test struct {
	DB         *gormx.DB
	HttpClient *httpx.Client
	Redis      *redisx.Store
}

var TestApiSet = wire.NewSet(wire.Struct(new(Test), "*"))

func (t *Test) TestApi(ctx *ws_context.WSContext) {
	fmt.Println(t.DB, t.HttpClient, t.Redis)
}
