package context

import "context"

type CronContext struct {
	context.Context
	Name string
	Spec string
	Job  interface{}
}
