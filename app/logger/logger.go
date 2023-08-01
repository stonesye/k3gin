package logger

import (
	"context"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"k3gin/app/config"
	"k3gin/app/contextx"
	"os"
	"path/filepath"
	"time"
)

// 设置下别名
type Entry = logrus.Entry
type Level = logrus.Level

// 1:fatal 2:error,3:warn,4:info,5:debug,6:trace
func SetLevel(level Level) {
	logrus.SetLevel(level)
}

func SetFormatter(format string) {
	switch format {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})
	}
}

// 作为logrus中的KEY
const (
	TraceIDKey = "trace_id"
	UserIDKey  = "user_id"
	TagKey     = "tag"
	StackKey   = "stack"
)

// WithContext 从ctx中把数据拿出来，封装成logrus的Fields, 这样后面再用logrus记录日志的时候，就会自动带上封装后的数据
func WithContext(ctx context.Context) *Entry {
	fields := logrus.Fields{}

	if v, e := contextx.FromTarceID(ctx); v != "" && e == true {
		fields[TraceIDKey] = v
	}

	if v := contextx.FromUserID(ctx); v != 0 {
		fields[UserIDKey] = v
	}

	if v, e := contextx.FromTag(ctx); v != "" && e == true {
		fields[TagKey] = v
	}

	if v := contextx.FromStack(ctx); v != nil {
		fields[StackKey] = fmt.Sprintf("%+v", v)
	}

	// 把上下文附带的给到logrus，以防万一后面要用
	return logrus.WithContext(ctx).WithFields(fields)
}

// InitLogger 初始化日志
func InitLogger() (func(), error) {
	c := config.C.Log

	SetLevel(Level(c.Level)) // 设置日志等级
	SetFormatter(c.Format)   // 设置日志格式

	var file *rotatelogs.RotateLogs

	if c.Output != "" { // 设置日志输出类型
		switch c.Output {
		case "stdout":
			logrus.SetOutput(os.Stdout)
		case "stderr":
			logrus.SetOutput(os.Stderr)
		case "file": // 当日志输出为文件的时候，借助rotatelogs包来完成日志分割
			if name := c.OutputFile; name != "" {
				_ = os.MkdirAll(filepath.Dir(name), 0777) // 创建日志存储目录

				// 利用rotatelogs包做日志分割
				file, err := rotatelogs.New(
					name+".%Y-%m-%d",
					rotatelogs.WithLinkName(name), // 日志文件地址
					rotatelogs.WithRotationTime(time.Duration(c.RotationTime)*time.Hour), // 日志轮训周期 一个日志文件存储多长时间
					rotatelogs.WithRotationCount(uint(c.RotationCount)))                  // 日志轮询数量, 一共存储多少个日志文件
				if err != nil {
					return nil, err
				}

				logrus.SetOutput(file)
			}
		}
	}

	return func() {
		if file != nil {
			_ = file.Close()
		}
	}, nil
}
