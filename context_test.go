package cronx

import (
	"context"
	"reflect"
	"testing"
)

func TestGetJobMetadata(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name  string
		args  args
		want  JobMetadata
		want1 bool
	}{
		{
			name:  "Nil",
			args:  args{},
			want:  JobMetadata{},
			want1: false,
		},
		{
			name: "Broken type",
			args: args{
				ctx: context.WithValue(context.Background(), CtxKeyJobMetadata, "this is string"),
			},
			want:  JobMetadata{},
			want1: false,
		},
		{
			name: "Exists",
			args: args{
				ctx: context.WithValue(context.Background(), CtxKeyJobMetadata, JobMetadata{
					EntryID:    1,
					Wave:       2,
					TotalWave:  3,
					IsLastWave: true,
				}),
			},
			want: JobMetadata{
				EntryID:    1,
				Wave:       2,
				TotalWave:  3,
				IsLastWave: true,
			},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetJobMetadata(tt.args.ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetJobMetadata() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetJobMetadata() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSetJobMetadata(t *testing.T) {
	type args struct {
		ctx  context.Context
		meta JobMetadata
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "Nil",
			args: args{
				ctx: nil,
				meta: JobMetadata{
					EntryID:    1,
					Wave:       2,
					TotalWave:  3,
					IsLastWave: true,
				},
			},
			want: context.WithValue(context.Background(), CtxKeyJobMetadata, JobMetadata{
				EntryID:    1,
				Wave:       2,
				TotalWave:  3,
				IsLastWave: true,
			}),
		},
		{
			name: "Exists",
			args: args{
				ctx: context.Background(),
				meta: JobMetadata{
					EntryID:    1,
					Wave:       2,
					TotalWave:  3,
					IsLastWave: true,
				},
			},
			want: context.WithValue(context.Background(), CtxKeyJobMetadata, JobMetadata{
				EntryID:    1,
				Wave:       2,
				TotalWave:  3,
				IsLastWave: true,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetJobMetadata(tt.args.ctx, tt.args.meta); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetJobMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}
