package main

import (
	"context"
	"errors"

	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/cronx/interceptor"
	"github.com/rizalgowandy/cronx/storage"
	"github.com/rizalgowandy/gdk/pkg/converter"
	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/fn"
	"github.com/rizalgowandy/gdk/pkg/logx"
	"github.com/rizalgowandy/gdk/pkg/storage/database"
)

type sendEmail struct{}

func (s sendEmail) Run(ctx context.Context) error {
	logx.INF(ctx, logx.KV{"job": fn.Name()}, "every 5 sec send reminder emails")
	return nil
}

type payBill struct{}

func (p payBill) Run(ctx context.Context) error {
	logx.INF(ctx, logx.KV{"job": fn.Name()}, "every 1 min pay bill")
	return nil
}

type alwaysError struct{}

func (a alwaysError) Run(ctx context.Context) error {
	err := errorx.E("some super long error message that come from executing the process")
	logx.ERR(ctx, err, "every 30 sec error")
	return err
}

type subscription struct{}

func (subscription) Run(ctx context.Context) error {
	md, ok := cronx.GetJobMetadata(ctx)
	if !ok {
		return errors.New("cannot job metadata")
	}
	logx.INF(ctx, logx.KV{"job": fn.Name(), "metadata": md}, "subscription is running")
	return nil
}

func callMe(ctx context.Context) error {
	logx.INF(ctx, logx.KV{"job": fn.Name()}, "call me every now and then")
	return nil
}

func main() {
	ctx := logx.NewContext()

	// Create database connection.
	db, err := database.NewPGXClient(ctx, &database.PostgreConfiguration{
		Address:               "user=unicorn_user password=magical_password dbname=example host=127.0.0.1 port=5432 sslmode=disable",
		MinConnection:         8,
		MaxConnection:         16,
		MaxConnectionLifetime: 3600,
		MaxConnectionIdleTime: 60,
	})
	if err != nil {
		logx.FTL(ctx, errorx.E(err), "new db client creation must success")
	}

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
		cronx.WithStorage(storage.NewPostgreClient(db)),
	)
	defer manager.Stop()

	// Register jobs.
	RegisterJobs(ctx, manager)

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

func RegisterJobs(ctx context.Context, manager *cronx.Manager) {
	// Struct name will become the name for the current job.
	if err := manager.Schedule("@every 5s", sendEmail{}); err != nil {
		// create log and send alert we fail to register job.
		logx.ERR(ctx, errorx.E(err), "register sendEmail must success")
	}

	// Create some jobs with the same struct.
	// Duplication is okay.
	for i := 0; i < 3; i++ {
		spec := "@every " + converter.String(i+1) + "m"
		if err := manager.Schedule(spec, payBill{}); err != nil {
			logx.ERR(ctx, errorx.E(err), "register payBill must success")
		}
	}

	// Create some jobs with broken spec.
	for i := 0; i < 2; i++ {
		spec := "broken spec " + converter.String(i+1)
		if err := manager.Schedule(spec, payBill{}); err != nil {
			logx.ERR(ctx, errorx.E(err), "register payBill must success")
		}
	}

	// Create a job with run that will always be error.
	if err := manager.Schedule("@every 30s", alwaysError{}); err != nil {
		logx.ERR(ctx, errorx.E(err), "register alwaysError must success")
	}

	// Create a job with v1 specification that includes seconds.
	if err := manager.Schedule("0 0 1 * * *", subscription{}); err != nil {
		logx.ERR(ctx, errorx.E(err), "register subscription must success")
	}

	// Create a job with multiple schedules
	if err := manager.Schedules("0 0 4 * * *#0 0 7 * * *#0 0 8 * * *", "#", subscription{}); err != nil {
		logx.ERR(ctx, errorx.E(err), "register subscription must success")
	}

	// Create a custom job with missing name.
	// A better approach is to used the cronx.ScheduleFunc so correct name can be shown.
	if err := manager.Schedule("0 */1 * * *", cronx.Func(func(context.Context) error {
		logx.INF(ctx, logx.KV{"job": "nameless job"}, "every 1h will be run")
		return nil
	})); err != nil {
		logx.ERR(ctx, errorx.E(err), "register job must success")
	}

	// Create a job with func instead of struct.
	if err := manager.ScheduleFunc("@every 10s", "callMe", callMe); err != nil {
		logx.ERR(ctx, errorx.E(err), "register callMe must success")
	}
	if err := manager.SchedulesFunc("0 0 9 * * *#0 0 10 * * *", "#", "callMe", callMe); err != nil {
		logx.ERR(ctx, errorx.E(err), "register callMe must success")
	}

	// Remove a job.
	const jobIDToBeRemoved = 2
	manager.Remove(jobIDToBeRemoved)

	// Get all current registered job.
	logx.INF(ctx, logx.KV{"entries": manager.GetEntries()}, "current jobs")
}
