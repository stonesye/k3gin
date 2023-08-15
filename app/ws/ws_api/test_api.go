package ws_api

import (
	"gorm.io/gorm"
	"k3gin/app/cache/redisx"
	"k3gin/app/httpx"
	"k3gin/app/ws/ws_context"
)

/**
测试websocket
*/

type Test struct {
	db         *gorm.DB
	httpClient *httpx.Client
	redis      *redisx.Store
}

func (t *Test) TestApi(ctx *ws_context.WSContext) {

}
