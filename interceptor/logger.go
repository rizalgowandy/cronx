package interceptor

import (
	"context"
	"fmt"
	"time"

	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/gdk/pkg/logx"
	"github.com/rizalgowandy/gdk/pkg/tags"
)

// Logger is a middleware that logs the current job start and finish.
func Logger() cronx.Interceptor {
	return func(ctx context.Context, job *cronx.Job, handler cronx.Handler) error {
		start := time.Now()
		err := handler(ctx, job)
		if err != nil {
			logx.ERR(ctx, err, job.Name)
			return err
		}

		logx.DBG(
			ctx,
			logx.KV{tags.Latency: time.Since(start).String()},
			fmt.Sprintf("operation cron %s success", job.Name),
		)
		return nil
	}
}
