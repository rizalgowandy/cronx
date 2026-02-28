package cronx

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rizalgowandy/cronx/storage"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type parserStub struct {
	err error
}

func (p parserStub) Parse(string) (cron.Schedule, error) {
	return nil, p.err
}

type storageStub struct{}

func (s storageStub) WriteHistory(context.Context, *storage.History) error {
	return nil
}

func (s storageStub) ReadHistories(context.Context, *storage.HistoryFilter) ([]storage.History, error) {
	return nil, nil
}

type alerterStub struct{}

func (a *alerterStub) NotifyHighLatency(context.Context, *Job, time.Time, time.Time, time.Duration, time.Duration) {
}

func TestWithParserOverride(t *testing.T) {
	t.Parallel()

	parserErr := errors.New("custom parser error")
	m := NewManager(
		WithAutoStartDisabled(),
		WithParser(parserStub{err: parserErr}),
	)

	err := m.Schedule("@every 1s", NewFuncJob("sample", func(context.Context) error {
		return nil
	}))
	require.Error(t, err)
	assert.ErrorIs(t, err, parserErr)
}

func TestWithInterceptor(t *testing.T) {
	t.Parallel()

	order := make([]string, 0, 2)
	m := NewManager(
		WithAutoStartDisabled(),
		WithInterceptor(func(ctx context.Context, job *Job, handler Handler) error {
			order = append(order, "interceptor")
			return handler(ctx, job)
		}),
	)

	err := m.interceptor(context.Background(), nil, func(ctx context.Context, job *Job) error {
		order = append(order, "handler")
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"interceptor", "handler"}, order)
}

func TestManagerOptionsAssignment(t *testing.T) {
	t.Parallel()

	loc := time.FixedZone("WIB", 7*60*60)
	storageClient := storageStub{}
	alerterClient := &alerterStub{}

	m := NewManager(
		WithAutoStartDisabled(),
		WithLocation(loc),
		WithLowPriorityDownJobs(),
		WithStorage(storageClient),
		WithAlerter(alerterClient),
	)

	assert.Equal(t, loc, m.location)
	assert.Equal(t, loc, m.createdTime.Location())
	assert.False(t, m.autoStart)
	assert.False(t, m.highPriorityDownJobs)
	assert.Same(t, alerterClient, m.alerter)
	assert.Equal(t, storageClient, m.storage)
}
