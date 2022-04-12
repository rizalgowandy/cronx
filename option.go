package cronx

import (
	"time"

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
func WithParser(p cron.Parser) Option {
	return func(m *Manager) {
		m.parser = p
	}
}

// WithChain specifies Job wrappers to apply to all jobs added to this cron.
func WithInterceptor(interceptors ...Interceptor) Option {
	return func(m *Manager) {
		m.interceptor = Chain(interceptors...)
	}
}
