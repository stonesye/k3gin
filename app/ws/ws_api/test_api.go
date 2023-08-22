package ws_api

import (
	"github.com/google/wire"
	"k3gin/app/cache/redisx"
	"k3gin/app/gormx"
	"k3gin/app/httpx"
	"k3gin/app/logger"
	"k3gin/app/util/trace"
	"k3gin/app/ws/ws_context"
	"k3gin/app/ws/ws_proxy"
	"time"
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
	ws, cleanFunc, err := ws_proxy.NewWebsocketConnect(ctx)
	if err != nil {
		return
	}
	defer cleanFunc()

	var hbTime = time.Second * time.Duration(10)

	hbTicker := time.NewTicker(hbTime)
	defer hbTicker.Stop()

	go func() {
		for {
			p := make([]byte, 1024)
			if _, err := ws.Read(p); err != nil {
				logger.WithFieldsFromContext(ctx).Errorf("Read client message err : %v", err)
				return
			}
			// 回写数据给客户端
			res := trace.NewTraceID() + "-" + string(p)
			if _, err := ws.Write([]byte(res)); err != nil {
				logger.WithFieldsFromContext(ctx).Errorf("write message err : %v", err)
				return
			}
			hbTicker.Reset(hbTime)
		}

	}()

	for {
		select {
		case <-hbTicker.C:
			ws.Write([]byte("链接长时间没有信息，已关闭."))
		}
	}

}
