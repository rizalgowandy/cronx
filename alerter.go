package cronx

import (
	"context"
	"fmt"
	"time"

	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/logx"
)

type AlerterItf interface {
	NotifyHighLatency(
		ctx context.Context,
		job *Job,
		prev, next time.Time,
		latency, maxLatency time.Duration,
	)
}

func NewAlerter() *Alerter {
	return &Alerter{}
}

type Alerter struct{}

func (a *Alerter) NotifyHighLatency(
	ctx context.Context,
	job *Job,
	prev, next time.Time,
	latency, maxLatency time.Duration,
) {
	logx.WRN(
		ctx,
		errorx.E(
			"current run has not finished before a new run for next schedule is started",
			errorx.Fields{
				"prev_schedule":   prev.String(),
				"next_schedule":   next.String(),
				"current_latency": latency.String(),
				"max_latency":     maxLatency.String(),
			},
		),
		fmt.Sprintf("Operation cron %s has high latency", job.Name),
	)
}
