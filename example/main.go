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
	"github.com/rizalgowandy/gdk/pkg/tags"
)

type alwaysDown struct{}

func (a alwaysDown) Run(_ context.Context) error {
	return nil
}

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
	err := errorx.E(
		"some super long error message that come from executing the process",
		errorx.Fields{tags.User: 1},
		errorx.MetricStatusExpectedErr,
		errorx.CodeNotFound,
	)
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

func callYou(ctx context.Context) error {
	logx.INF(ctx, logx.KV{"job": fn.Name()}, "call you every now and then")
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
	for i := 0; i < 2; i++ {
		spec := "@every " + converter.String(i+1) + "m"
		if err := manager.Schedule(spec, payBill{}); err != nil {
			logx.ERR(ctx, errorx.E(err), "register payBill must success")
		}
	}
	// Remove job with duplication.
	// See on the API or UI, entry_id = 2 will not be listed.
	const jobIDToBeRemoved = 2
	manager.Remove(jobIDToBeRemoved)

	// Create a job with broken spec.
	spec := "clearly a broken spec"
	if err := manager.Schedule(spec, alwaysDown{}); err != nil {
		logx.ERR(ctx, errorx.E(err), "register payBill must success")
	}

	// Create a job with run that will always be error.
	if err := manager.Schedule("@every 30s", alwaysError{}); err != nil {
		logx.ERR(ctx, errorx.E(err), "register alwaysError must success")
	}

	// Create a job with multiple schedules
	// The job uses v1 specification that includes seconds.
	if err := manager.Schedules("0 0 4 * * *#0 0 7 * * *#0 0 8 * * *", "#", subscription{}); err != nil {
		logx.ERR(ctx, errorx.E(err), "register subscription must success")
	}

	// Create a job with func instead of struct.
	if err := manager.ScheduleFunc("@every 10s", "callYou", callYou); err != nil {
		logx.ERR(ctx, errorx.E(err), "register callMe must success")
	}
	if err := manager.SchedulesFunc("0 0 9 * * *#0 0 10 * * *", "#", "callMe", callMe); err != nil {
		logx.ERR(ctx, errorx.E(err), "register callMe must success")
	}

	// Get all current registered job.
	logx.INF(ctx, logx.KV{"entries": manager.GetEntries()}, "current jobs")
}
