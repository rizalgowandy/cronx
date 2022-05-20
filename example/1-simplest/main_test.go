package main

import (
	"context"
	"testing"

	"github.com/rizalgowandy/cronx"
)

func Test_subscription_Run(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				in0: cronx.SetJobMetadata(context.Background(), cronx.JobMetadata{
					EntryID:    1,
					Wave:       2,
					TotalWave:  3,
					IsLastWave: true,
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			su := subscription{}
			if err := su.Run(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
