package main

import (
	"context"
	"errors"

	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/cronx/interceptor"
	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/fn"
	"github.com/rizalgowandy/gdk/pkg/logx"
)

type subscription struct{}

func (subscription) Run(ctx context.Context) error {
	md, ok := cronx.GetJobMetadata(ctx)
	if !ok {
		return errors.New("cannot job metadata")
	}
	logx.INF(ctx, logx.KV{"job": fn.Name(), "metadata": md}, "subscription is running")
	return nil
}

func main() {
	ctx := logx.NewContext()

	// Create middlewares.
	// The order is important.
	// The first one will be executed first.
	middlewares := cronx.Chain(
		interceptor.RequestID,
		interceptor.Recover(),
		interceptor.Logger(),
		interceptor.DefaultWorkerPool(),
	)

	// Create the manager with middleware.
	manager := cronx.NewManager(
		cronx.WithInterceptor(middlewares),
	)
	defer manager.Stop()

	// Create a job.
	if err := manager.Schedule("@every 1m", subscription{}); err != nil {
		logx.ERR(ctx, errorx.E(err), "register subscription must success")
	}

	// Get all current registered job.
	logx.INF(ctx, logx.KV{"entries": manager.GetEntries()}, "current jobs")

	// ===========================
	// Start Main Server
	// ===========================
	server, err := cronx.NewServer(manager, ":9001")
	if err != nil {
		logx.FTL(ctx, errorx.E(err), "new server creation must success")
		return
	}
	if err := server.ListenAndServe(); err != nil {
		logx.FTL(ctx, errorx.E(err), "server listen and server must success")
	}
}
