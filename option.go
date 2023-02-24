package cronx

import (
	"time"

	"github.com/rizalgowandy/cronx/storage"
	"github.com/robfig/cron/v3"
)

// Option represents a modification to the default behavior of the manager.
type Option func(*Manager)

// WithLocation overrides the timezone of the cron instance.
func WithLocation(loc *time.Location) Option {
	return func(m *Manager) {
		m.location = loc
	}
}

// WithParser overrides the parser used for interpreting job schedules.
func WithParser(p cron.ScheduleParser) Option {
	return func(m *Manager) {
		m.parser = p
	}
}

// WithInterceptor specifies Job wrappers to apply to all jobs added to this cron.
func WithInterceptor(interceptors ...Interceptor) Option {
	return func(m *Manager) {
		m.interceptor = Chain(interceptors...)
	}
}

// WithAutoStartDisabled prevent the cron job from actually running.
func WithAutoStartDisabled() Option {
	return func(m *Manager) {
		m.autoStart = false
	}
}

// WithLowPriorityDownJobs puts the down jobs at the bottom of the list.
func WithLowPriorityDownJobs() Option {
	return func(m *Manager) {
		m.highPriorityDownJobs = false
	}
}

// WithStorage determines the reader and writer for historical data.
func WithStorage(client storage.Client) Option {
	return func(m *Manager) {
		m.storage = client
	}
}

// WithAlerter determines the alerter used to send notification for high latency job run detected.
func WithAlerter(client AlerterItf) Option {
	return func(m *Manager) {
		m.alerter = client
	}
}
