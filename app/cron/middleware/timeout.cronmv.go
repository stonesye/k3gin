package middleware

import (
	"context"
	cronContext "k3gin/app/cron/context"
	"time"
)

func TimeoutCron(duration time.Duration) cronContext.HandleFunc {
	return func(c *cronContext.Context) {
		timeoutCtx, cancelFunc := context.WithTimeout(c, duration)
		defer cancelFunc()
		c.Context = timeoutCtx
		c.Next()
	}
}
