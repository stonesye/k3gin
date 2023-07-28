package trace

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

var (
	incrNum uint64
)

func NewTraceID() string {
	return fmt.Sprintf("traceID-%d-%s-%d",
		os.Getpid(),
		time.Now().Format("20060102150405.999"),
		atomic.AddUint64(&incrNum, 1))
}
