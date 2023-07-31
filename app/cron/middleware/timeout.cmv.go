package middleware

import (
	cronCtx "k3gin/app/cron/context"
	"time"
)

func TimeoutCron(duration time.Duration) cronCtx.HandleFunc {
	return func(cronCtx *cronCtx.CronContext) {
		cronCtx.Next()
	}
}
