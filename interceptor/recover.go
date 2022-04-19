package interceptor

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/gdk/pkg/env"
	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/logx"
	"github.com/rizalgowandy/gdk/pkg/netx"
	"github.com/rizalgowandy/gdk/pkg/stack"
	"github.com/rizalgowandy/gdk/pkg/tags"
)

// Recover is a middleware that recovers server from panic.
// Recover also dumps stack trace on panic occurrence.
func Recover() cronx.Interceptor {
	return func(ctx context.Context, job *cronx.Job, handler cronx.Handler) error {
		defer func() {
			if err := recover(); err != nil {
				fields := errorx.Fields{
					"stack": stack.ToArr(stack.Trim(debug.Stack())),
					"job":   job,
				}
				logx.ERR(ctx, errorx.E(err, fields), "recovered")
			}
		}()

		return handler(ctx, job)
	}
}

// RecoverWithAlert is a middleware that recovers server from panic.
// On panic occurrence, we will dump the stack trace and send alert to Slack.
func RecoverWithAlert(
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
		defer func() {
			if r := recover(); r != nil {
				stackTrace := stack.ToArr(stack.Trim(debug.Stack()))

				err = errorx.E(
					"there is a panic",
					errorx.Op(job.Name),
					errorx.CodeInternal,
					errorx.Fields{
						tags.StackTrace: stackTrace,
						tags.Panic:      r,
					})

				// Create a panic log.
				logx.ERR(
					ctx,
					err,
					"recovered from panic",
				)

				// Send slack alert.
				msg := fmt.Sprintf("[%s] *%v*", currentEnv, r)
				msg += fmt.Sprintf(" | `host: %s` | `app: %s-%s`", ipAddress, serviceName, mode)
				msg += fmt.Sprintf("\n*Job:* `%s`", job.Name)
				msg += fmt.Sprintf(
					"\n*Request ID:* `%s` _(search log file using this id)_",
					logx.GetRequestID(ctx),
				)
				msg += "\n*Stack Trace:*\n```"
				for k, v := range stackTrace {
					msg += v
					if k < len(stackTrace)-1 {
						msg += "\n"
					}
				}
				msg += "```"
				sc.Send(ctx, msg, isProduction)
			}
		}()

		return handler(ctx, job)
	}
}
