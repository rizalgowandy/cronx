package cronx

import (
	"context"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	type args struct {
		interceptors Interceptor
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
			got := NewManager(WithInterceptor(tt.args.interceptors))
			assert.NotNil(t, got)
		})
	}
}

func TestManager_Schedule(t *testing.T) {
	type args struct {
		spec string
		job  JobItf
		mock func() *Manager
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Broken spec",
			args: args{
				spec: "this is clearly not a spec",
				job:  Func(func(ctx context.Context) error { return nil }),
				mock: func() *Manager {
					manager := NewManager()
					return manager
				},
			},
			wantErr: true,
		},
		{
			name: "Success with descriptor",
			args: args{
				spec: "@every 5m",
				job:  Func(func(ctx context.Context) error { return nil }),
				mock: func() *Manager {
					manager := NewManager()
					return manager
				},
			},
			wantErr: false,
		},
		{
			name: "Success with v1",
			args: args{
				spec: "0 */30 * * * *",
				job:  Func(func(ctx context.Context) error { return nil }),
				mock: func() *Manager {
					manager := NewManager()
					return manager
				},
			},
			wantErr: false,
		},
		{
			name: "Success with v3",
			args: args{
				spec: "*/30 * * * *",
				job:  Func(func(ctx context.Context) error { return nil }),
				mock: func() *Manager {
					manager := NewManager()
					return manager
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := tt.args.mock()
			if err := manager.Schedule(tt.args.spec, tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("Schedule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_ScheduleFunc(t *testing.T) {
	type args struct {
		spec string
		name string
		cmd  func(ctx context.Context) error
		mock func() *Manager
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Broken spec",
			args: args{
				spec: "this is clearly not a spec",
				name: "nothing",
				cmd:  Func(func(ctx context.Context) error { return nil }),
				mock: func() *Manager {
					manager := NewManager()
					return manager
				},
			},
			wantErr: true,
		},
		{
			name: "Success with descriptor",
			args: args{
				spec: "@every 5m",
				name: "nothing",
				cmd:  Func(func(ctx context.Context) error { return nil }),
				mock: func() *Manager {
					manager := NewManager()
					return manager
				},
			},
			wantErr: false,
		},
		{
			name: "Success with v1",
			args: args{
				spec: "0 */30 * * * *",
				cmd:  Func(func(ctx context.Context) error { return nil }),
				mock: func() *Manager {
					manager := NewManager()
					return manager
				},
			},
			wantErr: false,
		},
		{
			name: "Success with v3",
			args: args{
				spec: "*/30 * * * *",
				name: "nothing",
				cmd:  Func(func(ctx context.Context) error { return nil }),
				mock: func() *Manager {
					manager := NewManager()
					return manager
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := tt.args.mock()
			if err := manager.ScheduleFunc(tt.args.spec, tt.args.name, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("ScheduleFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Schedules(t *testing.T) {
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
			manager := NewManager()
			if err := manager.Schedules(tt.args.spec, tt.args.separator, tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("Schedules() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_SchedulesFunc(t *testing.T) {
	type args struct {
		spec      string
		separator string
		name      string
		cmd       func(ctx context.Context) error
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
				cmd:       Func(func(ctx context.Context) error { return nil }),
			},
			wantErr: true,
		},
		{
			name: "Empty separator",
			args: args{
				spec:      "* 1 * * *",
				separator: "",
				cmd:       Func(func(ctx context.Context) error { return nil }),
			},
			wantErr: true,
		},
		{
			name: "Broken specification",
			args: args{
				spec:      "this is not specification#this is broken",
				separator: "#",
				cmd:       Func(func(ctx context.Context) error { return nil }),
			},
			wantErr: true,
		},
		{
			name: "Partial broken specification",
			args: args{
				spec:      "0 57 0 * * *#this is broken",
				separator: "#",
				cmd:       Func(func(ctx context.Context) error { return nil }),
			},
			wantErr: true,
		},
		{
			name: "Success with 1 waves",
			args: args{
				spec:      "0 57 0 * * *",
				separator: "#",
				cmd:       Func(func(ctx context.Context) error { return nil }),
			},
			wantErr: false,
		},
		{
			name: "Success with 3 waves",
			args: args{
				spec:      "0 57 0 * * *#0 18 16 * * *#0 7 1 * * *",
				separator: "#",
				cmd:       Func(func(ctx context.Context) error { return nil }),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager()
			if err := manager.SchedulesFunc(tt.args.spec, tt.args.separator, tt.args.name, tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("SchedulesFunc() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_Start(t *testing.T) {
	tests := []struct {
		name string
		mock func() *Manager
		want bool
	}{
		{
			name: "Success",
			mock: func() *Manager {
				manager := NewManager()
				_ = manager.Schedule(
					"@every 5m",
					Func(func(ctx context.Context) error { return nil }),
				)
				return manager
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := tt.mock()
			manager.Start()
		})
	}
}

func TestManager_Stop(t *testing.T) {
	tests := []struct {
		name string
		mock func() *Manager
		want bool
	}{
		{
			name: "Success",
			mock: func() *Manager {
				manager := NewManager()
				_ = manager.Schedule(
					"@every 5m",
					Func(func(ctx context.Context) error { return nil }),
				)
				return manager
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := tt.mock()
			manager.Stop()
		})
	}
}

func TestGetEntries(t *testing.T) {
	tests := []struct {
		name string
		mock func() *Manager
		want bool
	}{
		{
			name: "Success",
			mock: func() *Manager {
				manager := NewManager()
				_ = manager.Schedule(
					"@every 5m",
					Func(func(ctx context.Context) error { return nil }),
				)
				return manager
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := tt.mock()
			got := manager.GetEntries()
			if tt.want {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func TestManager_GetEntry(t *testing.T) {
	type args struct {
		id cron.EntryID
	}
	tests := []struct {
		name string
		args args
		mock func() *Manager
		want bool
	}{
		{
			name: "Success",
			mock: func() *Manager {
				manager := NewManager()
				_ = manager.Schedule(
					"@every 5m",
					Func(func(ctx context.Context) error { return nil }),
				)
				return manager
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := tt.mock()
			got := manager.GetEntry(tt.args.id)
			if tt.want {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func TestManager_Remove(t *testing.T) {
	type args struct {
		id cron.EntryID
	}
	tests := []struct {
		name string
		args args
		mock func() *Manager
		want bool
	}{
		{
			name: "Success",
			mock: func() *Manager {
				manager := NewManager()
				_ = manager.Schedule(
					"@every 5m",
					Func(func(ctx context.Context) error { return nil }),
				)
				return manager
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := tt.mock()
			manager.Remove(tt.args.id)
		})
	}
}

func TestManager_GetInfo(t *testing.T) {
	type fields struct {
		Commander   *cron.Cron
		Interceptor Interceptor
		Parser      cron.Parser
		DownJobs    []*Job
		Address     string
		Location    *time.Location
		CreatedTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "Success",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewManager()
			got := c.GetInfo()
			assert.NotNil(t, got)
		})
	}
}
