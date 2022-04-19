package interceptor

import (
	"context"
	"fmt"

	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/gdk/pkg/env"
	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/logx"
	"github.com/rizalgowandy/gdk/pkg/netx"
)

type SlackClientItf interface {
	Send(ctx context.Context, msg string, mention bool)
}

// NotifySlack is a middleware that send alert to slack on error.
func NotifySlack(
	serviceName string,
	mode string,
	sc SlackClientItf,
) cronx.Interceptor {
	var (
		currentEnv   = env.GetCurrent()
		isProduction = env.IsProduction()
		ipAddress    = netx.GetIPv4()
	)

	return func(ctx context.Context, job *cronx.Job, handler cronx.Handler) (err error) {
		err = handler(ctx, job)
		if err != nil {
			// No need to notify on expected error.
			if e, ok := err.(*errorx.Error); ok {
				if e.MetricStatus == errorx.MetricStatusExpectedErr {
					return err
				}
			}

			// Send slack alert.
			msg := fmt.Sprintf("[%s] *%v*", currentEnv, err.Error())
			msg += fmt.Sprintf(" | `host: %s` | `app: %s-%s`", ipAddress, serviceName, mode)
			msg += fmt.Sprintf("\n*Job:* `%s`", job.Name)
			msg += fmt.Sprintf(
				"\n*Request ID:* `%s` _(search log file using this id)_",
				logx.GetRequestID(ctx),
			)
			sc.Send(ctx, msg, isProduction)
		}
		return err
	}
}
