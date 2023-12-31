package job

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	croncontext "k3gin/app/cron/context"
	"k3gin/app/logger"
	"runtime"
	"time"
)

func TimeoutGlobalJob(duration time.Duration) func(ctx *croncontext.Context) {
	return func(ctx *croncontext.Context) {
		timeoutCtx, cancelFunc := context.WithTimeout(ctx.Context, duration)
		defer cancelFunc()
		ctx.Context = timeoutCtx
		ctx.Next()
	}
}

func RecoverGlobalJob() func(ctx *croncontext.Context) {
	return func(ctx *croncontext.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)
				logger.WithFieldsFromContext(ctx).Errorf("Cron painc err: %s", stack)
			}
		}()
		ctx.Next()
	}
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

// Job 包装一个对象，实现Run方法， 因为v3cron包里面的Job对象是需要实现Run函数的, 方便后续将任何待处理逻辑封装成Job，交给v3cron来处理
type Job struct {
	f func()
}

func (job *Job) Run() {
	job.f()
}

func NewJob(f func()) *Job {
	return &Job{f: f}
}
