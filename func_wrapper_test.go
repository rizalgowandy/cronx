package cronx

import (
	"context"
	"testing"
)

func TestFunc_Run(t *testing.T) {
	tests := []struct {
		name string
		r    Func
	}{
		{
			name: "Success",
			r:    Func(func(ctx context.Context) error { return nil }),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.r.Run(context.Background())
		})
	}
}
