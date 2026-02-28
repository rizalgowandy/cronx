package cronx

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChain(t *testing.T) {
	type args struct {
		interceptors []Interceptor
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
			args: args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Chain(tt.args.interceptors...)
			assert.NotNil(t, got)
		})
	}
}

func TestChainOrder(t *testing.T) {
	t.Parallel()

	calls := make([]string, 0, 5)
	handler := func(ctx context.Context, job *Job) error {
		calls = append(calls, "handler")
		return nil
	}

	chained := Chain(
		func(ctx context.Context, job *Job, next Handler) error {
			calls = append(calls, "first:before")
			err := next(ctx, job)
			calls = append(calls, "first:after")
			return err
		},
		func(ctx context.Context, job *Job, next Handler) error {
			calls = append(calls, "second:before")
			err := next(ctx, job)
			calls = append(calls, "second:after")
			return err
		},
	)

	err := chained(context.Background(), nil, handler)
	require.NoError(t, err)
	assert.Equal(t, []string{
		"first:before",
		"second:before",
		"handler",
		"second:after",
		"first:after",
	}, calls)
}

func TestChainErrorPropagation(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("handler failed")
	calls := make([]string, 0, 3)

	chained := Chain(
		func(ctx context.Context, job *Job, next Handler) error {
			calls = append(calls, "first")
			return next(ctx, job)
		},
		func(ctx context.Context, job *Job, next Handler) error {
			calls = append(calls, "second")
			return next(ctx, job)
		},
	)

	err := chained(context.Background(), nil, func(ctx context.Context, job *Job) error {
		calls = append(calls, "handler")
		return expectedErr
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, []string{"first", "second", "handler"}, calls)
}
