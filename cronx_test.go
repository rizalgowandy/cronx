package cronx

import (
	"context"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
)

func TestEvery(t *testing.T) {
	type args struct {
		duration time.Duration
		job      JobItf
		mock     func()
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Uninitialized",
			args: args{
				duration: 0,
				job:      nil,
				mock: func() {
					New()
					commandController.Commander = nil
				},
			},
		},
		{
			name: "Success",
			args: args{
				duration: 5 * time.Minute,
				job:      Func(func(ctx context.Context) error { return nil }),
				mock: func() {
					New()
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.mock()
			Every(tt.args.duration, tt.args.job)
		})
	}
}

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

func TestGetEntries(t *testing.T) {
	tests := []struct {
		name string
		mock func()
		want bool
	}{
		{
			name: "Uninitialized",
			mock: func() {
				New()
				commandController.Commander = nil
			},
		},
		{
			name: "Success",
			mock: func() {
				New()
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got := GetEntries()
			if tt.want {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func TestRemove(t *testing.T) {
	type args struct {
		id   cron.EntryID
		mock func()
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Uninitialized",
			args: args{
				id: 1,
				mock: func() {
					New()
					commandController.Commander = nil
				},
			},
		},
		{
			name: "Success",
			args: args{
				id: 1,
				mock: func() {
					New()
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.mock()
			Remove(tt.args.id)
		})
	}
}

func TestSchedule(t *testing.T) {
	type args struct {
		spec string
		job  JobItf
		mock func()
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Uninitialized",
			args: args{
				spec: "@every 5m",
				job:  Func(func(ctx context.Context) error { return nil }),
				mock: func() {
					New()
					commandController.Commander = nil
				},
			},
			wantErr: true,
		},
		{
			name: "Broken spec",
			args: args{
				spec: "this is clearly not a spec",
				job:  Func(func(ctx context.Context) error { return nil }),
				mock: func() {
					New()
				},
			},
			wantErr: true,
		},
		{
			name: "Success with descriptor",
			args: args{
				spec: "@every 5m",
				job:  Func(func(ctx context.Context) error { return nil }),
				mock: func() {
					New()
				},
			},
			wantErr: false,
		},
		{
			name: "Success with v1",
			args: args{
				spec: "0 */30 * * * *",
				job:  Func(func(ctx context.Context) error { return nil }),
				mock: func() {
					New()
				},
			},
			wantErr: false,
		},
		{
			name: "Success with v3",
			args: args{
				spec: "*/30 * * * *",
				job:  Func(func(ctx context.Context) error { return nil }),
				mock: func() {
					New()
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.mock()
			if err := Schedule(tt.args.spec, tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("Schedule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefault(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Default()
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
			args: args{
				config: Config{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			New()
		})
	}
}

func TestCustom(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Success",
			args: args{
				config: Config{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Custom(tt.args.config)
		})
	}
}

func TestStop(t *testing.T) {
	tests := []struct {
		name string
		mock func()
	}{
		{
			name: "Uninitialized",
			mock: func() {
				New()
				commandController.Commander = nil
			},
		},
		{
			name: "Success",
			mock: func() {
				New()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			Stop()
		})
	}
}

func TestSchedules(t *testing.T) {
	type args struct {
		spec      string
		separator string
		job       JobItf
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Empty specification",
			args: args{
				spec:      "",
				separator: "#",
				job:       Func(func(ctx context.Context) error { return nil }),
			},
			wantErr: true,
		},
		{
			name: "Empty separator",
			args: args{
				spec:      "* 1 * * *",
				separator: "",
				job:       Func(func(ctx context.Context) error { return nil }),
			},
			wantErr: true,
		},
		{
			name: "Broken specification",
			args: args{
				spec:      "this is not specification#this is broken",
				separator: "#",
				job:       Func(func(ctx context.Context) error { return nil }),
			},
			wantErr: true,
		},
		{
			name: "Partial broken specification",
			args: args{
				spec:      "0 57 0 * * *#this is broken",
				separator: "#",
				job:       Func(func(ctx context.Context) error { return nil }),
			},
			wantErr: true,
		},
		{
			name: "Success with 1 waves",
			args: args{
				spec:      "0 57 0 * * *",
				separator: "#",
				job:       Func(func(ctx context.Context) error { return nil }),
			},
			wantErr: false,
		},
		{
			name: "Success with 3 waves",
			args: args{
				spec:      "0 57 0 * * *#0 18 16 * * *#0 7 1 * * *",
				separator: "#",
				job:       Func(func(ctx context.Context) error { return nil }),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Schedules(tt.args.spec, tt.args.separator, tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("Schedules() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetStatusData(t *testing.T) {
	tests := []struct {
		name string
		mock func()
		want bool
	}{
		{
			name: "Uninitialized",
			mock: func() {
				commandController = nil
			},
		},
		{
			name: "Success without any job",
			mock: func() {
				New()
			},
			want: true,
		},
		{
			name: "Success",
			mock: func() {
				New()
				_ = Schedule("@every 5m", Func(func(ctx context.Context) error { return nil }))
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got := GetStatusData()
			if tt.want {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func TestGetStatusJSON(t *testing.T) {
	tests := []struct {
		name string
		mock func()
		want bool
	}{
		{
			name: "Uninitialized",
			mock: func() {
				commandController = nil
			},
		},
		{
			name: "Success without any job",
			mock: func() {
				New()
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got := GetStatusJSON()
			if tt.want {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func TestGetEntry(t *testing.T) {
	type args struct {
		id cron.EntryID
	}
	tests := []struct {
		name string
		args args
		mock func()
		want bool
	}{
		{
			name: "Uninitialized",
			mock: func() {
				New()
				commandController.Commander = nil
			},
		},
		{
			name: "Success",
			mock: func() {
				New()
				Every(time.Hour, Func(func(ctx context.Context) error {
					return nil
				}))
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got := GetEntry(tt.args.id)
			if tt.want {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func TestGetInfo(t *testing.T) {
	tests := []struct {
		name string
		mock func()
		want bool
	}{
		{
			name: "Uninitialized",
			mock: func() {
				commandController = nil
			},
		},
		{
			name: "Success without any job",
			mock: func() {
				New()
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got := GetInfo()
			if tt.want {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}
