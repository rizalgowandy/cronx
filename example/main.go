package main

import (
	"context"
	"errors"

	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/cronx/interceptor"
	"github.com/rizalgowandy/gdk/pkg/converter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type sendEmail struct{}

func (s sendEmail) Run(context.Context) error {
	log.WithLevel(zerolog.InfoLevel).
		Str("job", "sendEmail").
		Msg("every 5 sec send reminder emails")
	return nil
}

type payBill struct{}

func (p payBill) Run(context.Context) error {
	log.WithLevel(zerolog.InfoLevel).
		Str("job", "payBill").
		Msg("every 1 min pay bill")
	return nil
}

type alwaysError struct{}

func (a alwaysError) Run(context.Context) error {
	log.WithLevel(zerolog.InfoLevel).
		Str("job", "alwaysError").
		Msg("every 30 sec error")
	return errors.New("some super long error message that come from executing the process")
}

type everyJob struct{}

func (everyJob) Run(context.Context) error {
	log.WithLevel(zerolog.InfoLevel).
		Str("job", "everyJob").
		Msg("is running")
	return nil
}

type subscription struct{}

func (subscription) Run(ctx context.Context) error {
	md, ok := cronx.GetJobMetadata(ctx)
	if !ok {
		return errors.New("cannot job metadata")
	}

	log.WithLevel(zerolog.InfoLevel).
		Str("job", "subscription").
		Interface("metadata", md).
		Msg("is running")
	return nil
}

func main() {
	// Create middlewares.
	// The order is important.
	// The first one will be executed first.
	middlewares := cronx.Chain(
		interceptor.Recover(),
		interceptor.Logger(),
		interceptor.DefaultWorkerPool(),
	)

	// Create the manager with middleware.
	manager := cronx.NewManager(cronx.Config{}, middlewares)
	defer manager.Stop()

	// Register jobs.
	RegisterJobs(manager)

	// ===========================
	// Start Main Server
	// ===========================
	server, err := cronx.NewServer(manager, ":9001")
	if err != nil {
		log.WithLevel(zerolog.FatalLevel).
			Err(err).
			Msg("new server creation must success")
		return
	}
	if err := server.ListenAndServe(); err != nil {
		log.WithLevel(zerolog.FatalLevel).
			Err(err).
			Msg("server listen and server must success")
	}
}

func RegisterJobs(manager *cronx.Manager) {
	// Struct name will become the name for the current job.
	if err := manager.Schedule("@every 5s", sendEmail{}); err != nil {
		// create log and send alert we fail to register job.
		log.WithLevel(zerolog.ErrorLevel).
			Err(err).
			Msg("register sendEmail must success")
	}

	// Create some jobs with the same struct.
	// Duplication is okay.
	for i := 0; i < 3; i++ {
		spec := "@every " + converter.String(i+1) + "m"
		if err := manager.Schedule(spec, payBill{}); err != nil {
			log.WithLevel(zerolog.ErrorLevel).
				Err(err).
				Msg("register payBill must success")
		}
	}

	// Create some jobs with broken spec.
	for i := 0; i < 3; i++ {
		spec := "broken spec " + converter.String(i+1)
		if err := manager.Schedule(spec, payBill{}); err != nil {
			log.WithLevel(zerolog.ErrorLevel).
				Err(err).
				Msg("register payBill must success")
		}
	}

	// Create a job with run that will always be error.
	if err := manager.Schedule("@every 30s", alwaysError{}); err != nil {
		log.WithLevel(zerolog.ErrorLevel).
			Err(err).
			Msg("register alwaysError must success")
	}

	// Create a custom job with missing name.
	if err := manager.Schedule("0 */1 * * *", cronx.Func(func(context.Context) error {
		log.WithLevel(zerolog.InfoLevel).
			Str("job", "nameless job").
			Msg("every 1h will be run")
		return nil
	})); err != nil {
		log.WithLevel(zerolog.ErrorLevel).
			Err(err).
			Msg("register job must success")
	}

	// Create a job with v1 specification that includes seconds.
	if err := manager.Schedule("0 0 1 * * *", subscription{}); err != nil {
		log.WithLevel(zerolog.ErrorLevel).
			Err(err).
			Msg("register subscription must success")
	}

	// Create a job with multiple schedules
	if err := manager.Schedules("0 0 4 * * *#0 0 7 * * *#0 0 11 * * *", "#", subscription{}); err != nil {
		log.WithLevel(zerolog.ErrorLevel).
			Err(err).
			Msg("register subscription must success")
	}

	// Remove a job.
	const jobIDToBeRemoved = 2
	manager.Remove(jobIDToBeRemoved)

	// Get all current registered job.
	log.WithLevel(zerolog.InfoLevel).
		Interface("entries", manager.GetEntries()).
		Msg("current jobs")
}
