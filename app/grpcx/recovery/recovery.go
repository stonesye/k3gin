package recovery

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(opts ...func(*options)) grpc.UnaryServerInterceptor {
	o := initOptions(opts...)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = packRecoverMessage(ctx, r, o.recoveryHandlerFunc)
			}
		}()
		resp, err = handler(ctx, req)
		return resp, err
	}

}

func StreamServerInterceptor(opt ...func(*options)) grpc.StreamServerInterceptor {
	o := initOptions(opt...)

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {

		defer func() {
			if r := recover(); r != nil {
				err = packRecoverMessage(ss.Context(), r, o.recoveryHandlerFunc)
			}
		}()

		err = handler(srv, ss)
		return err
	}
}

func packRecoverMessage(ctx context.Context, p interface{}, f func(context.Context, interface{}) error) error {

	if f == nil {
		return status.Errorf(codes.Internal, "%v", p)
	}

	return f(ctx, p)
}
