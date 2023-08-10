package recovery

import (
	"context"
)

type options struct {
	recoveryHandlerFunc func(context.Context, interface{}) error
}

func WithRecoveryHandler(recoveryHandler func(interface{}) error) func(*options) {
	return func(o *options) {
		o.recoveryHandlerFunc = func(ctx context.Context, i interface{}) error {
			return recoveryHandler(i)
		}
	}
}

func WithRecoverHandlerContext(recoveryHandlerContext func(context.Context, interface{}) error) func(*options) {
	return func(o *options) {
		o.recoveryHandlerFunc = recoveryHandlerContext
	}
}

func initOptions(opts ...func(*options)) *options {
	var o = options{recoveryHandlerFunc: nil}

	for _, opt := range opts {
		opt(&o)
	}

	return &o

}
