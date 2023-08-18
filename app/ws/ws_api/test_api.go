package ws_api

import (
	"github.com/google/wire"
	"k3gin/app/cache/redisx"
	"k3gin/app/gormx"
	"k3gin/app/httpx"
	"k3gin/app/logger"
	"k3gin/app/ws/ws_context"
	"k3gin/app/ws/ws_router"
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
	ws, cleanFunc, err := ws_router.NewWebsocketConnect(ctx)
	if err != nil {
		return
	}

	defer cleanFunc()

	var p = make([]byte, 1024)

	for {
		n, err := ws.Read(p)
		if err != nil {
			logger.WithFieldsFromWSContext(ctx).Errorf("Read message err[%v] : %v", n, err)
			break
		}

		n, err = ws.Write(p)
		if err != nil {
			logger.WithFieldsFromWSContext(ctx).Errorf("Recv message err[%v] : %v", n, err)
			break
		}
	}

}
