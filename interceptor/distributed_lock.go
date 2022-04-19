package interceptor

import (
	"context"
	"fmt"

	"github.com/go-redsync/redsync/v4"
	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/gdk/pkg/env"
	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/logx"
	"github.com/rizalgowandy/gdk/pkg/netx"
	"github.com/rizalgowandy/gdk/pkg/tags"
)

type SlackClientItf interface {
	Send(ctx context.Context, msg string, mention bool)
}

type DistributedLockItf interface {
	Mutex(name string) *redsync.Mutex
}

// DistributedLock is a middleware that prevents a process executed at the same time across servers.
func DistributedLock(
	serviceName string,
	mode string,
	dl DistributedLockItf,
	sc SlackClientItf,
) cronx.Interceptor {
	var (
		currentEnv   = env.GetCurrent()
		isProduction = env.IsProduction()
		ipAddress    = netx.GetIPv4()
	)

	return func(ctx context.Context, job *cronx.Job, handler cronx.Handler) error {
		// Create lock.
		mutex := dl.Mutex(job.Name)
		if err := mutex.LockContext(ctx); err != nil {
			logx.ERR(
				ctx,
				errorx.E(err, errorx.Op(job.Name), errorx.CodeInternal, errorx.Fields{
					tags.Address: ipAddress,
				}),
				"distributed lock: cannot gain lock for current process",
			)
			return err
		}

		jobErr := handler(ctx, job)

		// Unlock.
		if _, err := mutex.UnlockContext(ctx); err != nil {
			logx.ERR(
				ctx,
				errorx.E(err, errorx.Op(job.Name), errorx.CodeInternal, errorx.Fields{
					tags.Address: ipAddress,
				}),
				"distributed lock: cannot unlock process",
			)

			// Send slack alert when failed to unlock.
			msg := fmt.Sprintf("[%s] *%v*", currentEnv, err)
			msg += fmt.Sprintf(" | `host: %s` | `app: %s-%s`", ipAddress, serviceName, mode)
			msg += fmt.Sprintf("\n*Job:* `%s`", job.Name)
			msg += fmt.Sprintf(
				"\n*Request ID:* `%s` _(search log file using this id)_",
				logx.GetRequestID(ctx),
			)
			sc.Send(ctx, msg, isProduction)
		}

		return jobErr
	}
}
