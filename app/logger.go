package app

import (
	"context"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"k3gin/app/config"
	"os"
	"path/filepath"
	"time"
)

// 设置下别名

type Logger = logrus.Logger
type Entry = logrus.Entry
type Hook = logrus.Hook
type Level = logrus.Level

// 各日志级别
const (
	PANIC Level = iota
	FATAL
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
)

func SetLevel(level Level) {
	logrus.SetLevel(level)
}

func SetFormatter(format string) {
	switch format {
	case "json":
		logrus.SetFormatter(new(logrus.JSONFormatter))
	default:
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}
}

// 定义穿插在日志和Context中的标记的Key
const (
	TraceIDKey  = "trace_id"
	UserIDKey   = "user_id"
	UserNameKey = "user_name"
	TagKey      = "tag"
	StackKey    = "stack"
)

// NewTraceIDContext 创建一个存放TraceID Context
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// FromTraceIDContext 获取context中的tarce_id的值
func FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(TraceIDKey)
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewUserIDContext(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

func FromUserIDContext(ctx context.Context) uint64 {
	v := ctx.Value(UserIDKey)
	if v != nil {
		if s, ok := v.(uint64); ok {
			return s
		}
	}

	return 0
}

func NewUserNameContext(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, UserNameKey, userName)
}

func FromUserNameContext(ctx context.Context) string {
	v := ctx.Value(UserNameKey)
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewTagContext(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, TagKey, tag)
}

func FromTagContext(ctx context.Context) string {
	v := ctx.Value(TagKey)
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewStackContext(ctx context.Context, stack error) context.Context {
	return context.WithValue(ctx, StackKey, stack)
}

func FromStackContext(ctx context.Context) error {
	v := ctx.Value(StackKey)
	if v != nil {
		if s, ok := v.(error); ok {
			return s
		}
	}
	return nil
}

// WithContext 从ctx中把数据拿出来，封装成logrus的Fields, 这样后面再用logrus记录日志的时候，就会自动带上封装后的数据
func WithContext(ctx context.Context) *Entry {
	fields := logrus.Fields{}

	if v := FromTraceIDContext(ctx); v != "" {
		fields[TraceIDKey] = v
	}

	if v := FromUserIDContext(ctx); v != 0 {
		fields[UserIDKey] = v
	}

	if v := FromUserNameContext(ctx); v != "" {
		fields[UserNameKey] = v
	}

	if v := FromTagContext(ctx); v != "" {
		fields[TagKey] = v
	}

	if v := FromStackContext(ctx); v != nil {
		fields[StackKey] = fmt.Sprintf("%+v", v)
	}

	// 把上下文附带的给到logrus，以防万一后面要用
	return logrus.WithContext(ctx).WithFields(fields)
}

// InitLogger 初始化日志
func InitLogger() (func(), error) {
	c := config.C.Log

	SetLevel(Level(c.Level))
	SetFormatter(c.Format)

	var file *rotatelogs.RotateLogs

	if c.Output != "" {
		switch c.Output {
		case "stdout":
			logrus.SetOutput(os.Stdout)
		case "stderr":
			logrus.SetOutput(os.Stderr)
		case "file":
			if name := c.OutputFile; name != "" {
				_ = os.MkdirAll(filepath.Dir(name), 0777)

				// 利用rotatelogs包做日志分割
				file, err := rotatelogs.New(
					name+".%Y-%m-%d",
					rotatelogs.WithLinkName(name),
					rotatelogs.WithRotationTime(time.Duration(c.RotationTime)*time.Hour),
					rotatelogs.WithRotationCount(uint(c.RotationCount)))
				if err != nil {
					return nil, err
				}

				logrus.SetOutput(file)
			}
		}
	}

	return func() {
		if file != nil {
			file.Close()
		}
	}, nil
}
