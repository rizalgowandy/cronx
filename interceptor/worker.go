package interceptor

import (
	"context"

	"github.com/rizalgowandy/cronx"
)

// Default configuration.
var defaultWorkerPoolSize = 1000

// WorkerPool is a middleware that limit total cron that can run a time.
// Program is running on a server with finite amount of resources such as CPU and RAM.
// By limiting the total number of jobs that can be run the same time,
// we protect the server from overloading.
func WorkerPool(size int) cronx.Interceptor {
	if size <= 0 {
		size = defaultWorkerPoolSize
	}
	pool := make(chan struct{}, size)

	return func(ctx context.Context, job *cronx.Job, handler cronx.Handler) error {
		// Wait for worker to be available.
		pool <- struct{}{}

		// Release the worker.
		defer func() {
			<-pool
		}()

		return handler(ctx, job)
	}
}

// DefaultWorkerPool returns a WorkerPool middleware with default configuration.
func DefaultWorkerPool() cronx.Interceptor {
	return WorkerPool(defaultWorkerPoolSize)
}
