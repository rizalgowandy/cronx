package cronx

import "context"

type (
	// Handler is the handler definition to run a job.
	Handler func(ctx context.Context, job *Job) error

	// Interceptor is the middleware that will be executed before the current handler.
	Interceptor func(ctx context.Context, job *Job, handler Handler) error
)

// Chain returns a single interceptor from multiple interceptors.
func Chain(interceptors ...Interceptor) Interceptor {
	n := len(interceptors)

	return func(ctx context.Context, job *Job, handler Handler) error {
		chainer := func(currentInter Interceptor, currentHandler Handler) Handler {
			return func(currentCtx context.Context, currentJob *Job) error {
				return currentInter(currentCtx, currentJob, currentHandler)
			}
		}

		chainedHandler := handler
		for i := n - 1; i >= 0; i-- {
			chainedHandler = chainer(interceptors[i], chainedHandler)
		}

		return chainedHandler(ctx, job)
	}
}
