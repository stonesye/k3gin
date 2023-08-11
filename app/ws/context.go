package ws

import (
	"k3gin/app/cron/context"
	"net/http"
)

// Context  websocket context
type Context struct {
	context.Context
	Request  *http.Request
	Response http.ResponseWriter
	Params   Params
}

type Params []Param

type Param struct {
	Key   string
	Value string
}
