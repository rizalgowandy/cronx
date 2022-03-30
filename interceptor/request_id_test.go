package interceptor

import (
	"context"
	"testing"

	"github.com/rizalgowandy/cronx"
	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
)

func TestRequestID(t *testing.T) {
	type args struct {
		ctx     context.Context
		job     *cronx.Job
		handler cronx.Handler
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error",
			args: args{
				ctx: context.Background(),
				job: &cronx.Job{},
				handler: func(ctx context.Context, job *cronx.Job) error {
					return errorx.E("error")
				},
			},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				job: &cronx.Job{},
				handler: func(ctx context.Context, job *cronx.Job) error {
					return nil
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RequestID(tt.args.ctx, tt.args.job, tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("RequestID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
