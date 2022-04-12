package interceptor

import (
	"context"
	"runtime/debug"

	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/logx"
	"github.com/rizalgowandy/gdk/pkg/stack"
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
