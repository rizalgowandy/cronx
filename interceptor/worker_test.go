package interceptor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool(t *testing.T) {
	type args struct {
		poolSize int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success with default configuration",
			args: args{},
		},
		{
			name: "Success",
			args: args{
				poolSize: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WorkerPool(tt.args.poolSize)
			assert.NotNil(t, got)
		})
	}
}

func TestDefaultWorkerPool(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultWorkerPool()
			assert.NotNil(t, got)
		})
	}
}
