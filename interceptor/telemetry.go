package interceptor

import (
	"context"
	"time"

	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/tags"
)

const (
	MetricCronPerformance = "cron_performance"
)

type TelemetryClientItf interface {
	Histogram(metricName string, value float64, tags map[string]string)
}

// Telemetry is a middleware that push the latency of a process.
func Telemetry(metric TelemetryClientItf) cronx.Interceptor {
	const dividerNStoMS = 1e6

	return func(ctx context.Context, job *cronx.Job, handler cronx.Handler) error {
		start := time.Now()
		err := handler(ctx, job)

		// Publish current metric.
		go func() {
			metric.Histogram(
				MetricCronPerformance,
				float64(time.Since(start).Nanoseconds()/dividerNStoMS),
				DefaultKV(job.Name, err),
			)
		}()

		return err
	}
}

// DefaultKV returns a standardized service tracking.
func DefaultKV(method string, err error) map[string]string {
	trackingTags := map[string]string{
		tags.Method:       method,
		tags.MetricStatus: string(errorx.MetricStatusSuccess),
	}

	// Set status as error.
	if err != nil {
		trackingTags[tags.Error] = err.Error()
		trackingTags[tags.MetricStatus] = string(errorx.MetricStatusErr)

		// If error is our custom error, add additional tags.
		if e, ok := err.(*errorx.Error); ok {
			if e.MetricStatus != "" {
				trackingTags[tags.MetricStatus] = string(e.MetricStatus)
			}
			if e.Code != "" {
				trackingTags[tags.Code] = string(e.Code)
			}
			if e.Line != "" {
				trackingTags[tags.ErrorLine] = string(e.Line)
			}
			if e.Message != "" {
				trackingTags[tags.Message] = string(e.Message)
			}
			if len(e.OpTraces) > 0 {
				trackingTags[tags.Ops] = string(e.OpTraces[len(e.OpTraces)-1])
			}
		}
	}

	return trackingTags
}
