package interceptor

import (
	"context"

	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/gdk/pkg/logx"
)

// RequestID is a middleware that inject request id to the context if it doesn't exists.
func RequestID(
	ctx context.Context,
	job *cronx.Job,
	handler cronx.Handler,
) error {
	return handler(logx.ContextWithRequestID(ctx), job)
}
