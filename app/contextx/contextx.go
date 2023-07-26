package contextx

import "context"

type (
	// 作为context中的KEY
	traceCtx  struct{} // 贯穿每个Request请求追踪ID
	userIDCtx struct{} // 贯穿每个Request请求UserID
	tagCtx    struct{} // 贯穿每个Request请求记录文件名
	stackCtx  struct{} // 贯穿每个Request请求记录堆栈错误
)

func NewStack(ctx context.Context, stack error) context.Context {
	return context.WithValue(ctx, stackCtx{}, stack)
}

func FromStack(ctx context.Context) error {
	v := ctx.Value(stackCtx{})
	if s, ok := v.(error); ok {
		return s
	}
	return nil
}

func NewTag(ctx context.Context, tagName string) context.Context {
	return context.WithValue(ctx, tagCtx{}, tagName)
}

func FromTag(ctx context.Context) (string, bool) {
	v := ctx.Value(tagCtx{})

	if s, ok := v.(string); ok {
		return s, s != ""
	}
	return "", false
}

func NewUserID(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, userIDCtx{}, userID)
}

func FromUserID(ctx context.Context) uint64 {
	v := ctx.Value(userIDCtx{})
	if v != nil {
		if s, ok := v.(uint64); ok {
			return s
		}
	}
	return 0
}

func NewTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceCtx{}, traceID)
}

func FromTarceID(ctx context.Context) (string, bool) {
	v := ctx.Value(traceCtx{})
	if s, ok := v.(string); ok {
		return s, s != ""
	}
	return "", false
}
