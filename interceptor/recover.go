package interceptor

import (
	"context"
	"runtime/debug"

	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/gdk/pkg/stack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Recover is a middleware that recovers server from panic.
// Recover also dumps stack trace on panic occurrence.
func Recover() cronx.Interceptor {
	return func(ctx context.Context, job *cronx.Job, handler cronx.Handler) error {
		defer func() {
			if err := recover(); err != nil {
				log.WithLevel(zerolog.PanicLevel).
					Interface("err", err).
					Interface("stack", stack.ToArr(stack.Trim(debug.Stack()))).
					Interface("job", job).
					Msg("recovered")
			}
		}()

		return handler(ctx, job)
	}
}
