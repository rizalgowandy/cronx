package cronx

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	type args struct {
		config       Config
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
			got := NewManager(tt.args.config, tt.args.interceptors)
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
					manager := NewManager(Config{})
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
					manager := NewManager(Config{})
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
					manager := NewManager(Config{})
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
					manager := NewManager(Config{})
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
			manager := NewManager(Config{})
			if err := manager.Schedules(tt.args.spec, tt.args.separator, tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("Schedules() error = %v, wantErr %v", err, tt.wantErr)
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
				manager := NewManager(Config{})
				_ = manager.Schedule("@every 5m", Func(func(ctx context.Context) error { return nil }))
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
				manager := NewManager(Config{})
				_ = manager.Schedule("@every 5m", Func(func(ctx context.Context) error { return nil }))
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
				manager := NewManager(Config{})
				_ = manager.Schedule("@every 5m", Func(func(ctx context.Context) error { return nil }))
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
				manager := NewManager(Config{})
				_ = manager.Schedule("@every 5m", Func(func(ctx context.Context) error { return nil }))
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
				manager := NewManager(Config{})
				_ = manager.Schedule("@every 5m", Func(func(ctx context.Context) error { return nil }))
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
		Commander        *cron.Cron
		Interceptor      Interceptor
		Parser           cron.Parser
		UnregisteredJobs []*Job
		Address          string
		Location         *time.Location
		CreatedTime      time.Time
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
			c := &Manager{
				Commander:        tt.fields.Commander,
				Interceptor:      tt.fields.Interceptor,
				Parser:           tt.fields.Parser,
				UnregisteredJobs: tt.fields.UnregisteredJobs,
				Location:         tt.fields.Location,
				CreatedTime:      tt.fields.CreatedTime,
			}
			got := c.GetInfo()
			assert.NotNil(t, got)
		})
	}
}

func TestManager_GetStatusData(t *testing.T) {
	type fields struct {
		Commander        *cron.Cron
		Interceptor      Interceptor
		Parser           cron.Parser
		UnregisteredJobs []*Job
	}
	tests := []struct {
		name   string
		fields fields
		want   []StatusData
	}{
		{
			name:   "Cron is nil",
			fields: fields{},
			want:   nil,
		},
		{
			name: "Success",
			fields: fields{
				Commander:   cron.New(),
				Interceptor: nil,
				Parser:      cron.Parser{},
				UnregisteredJobs: []*Job{
					{
						Name:    "Cron 1",
						Status:  "DOWN",
						Latency: "",
						Error:   "",
						inner:   nil,
						status:  statusDown,
						running: sync.Mutex{},
					},
				},
			},
			want: []StatusData{
				{
					ID: 0,
					Job: &Job{
						Name:    "Cron 1",
						Status:  "DOWN",
						Latency: "",
						Error:   "",
						inner:   nil,
						status:  statusDown,
						running: sync.Mutex{},
					},
					Next: time.Time{},
					Prev: time.Time{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Manager{
				Commander:        tt.fields.Commander,
				Interceptor:      tt.fields.Interceptor,
				Parser:           tt.fields.Parser,
				UnregisteredJobs: tt.fields.UnregisteredJobs,
			}
			if got := c.GetStatusData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStatusData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_GetStatusJSON(t *testing.T) {
	type fields struct {
		Commander        *cron.Cron
		Interceptor      Interceptor
		Parser           cron.Parser
		UnregisteredJobs []*Job
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "Cron is nil",
			fields: fields{},
		},
		{
			name: "Success",
			fields: fields{
				Commander:   cron.New(),
				Interceptor: nil,
				Parser:      cron.Parser{},
				UnregisteredJobs: []*Job{
					{
						Name:    "Cron 1",
						Status:  "DOWN",
						Latency: "",
						Error:   "",
						inner:   nil,
						status:  statusDown,
						running: sync.Mutex{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Manager{
				Commander:        tt.fields.Commander,
				Interceptor:      tt.fields.Interceptor,
				Parser:           tt.fields.Parser,
				UnregisteredJobs: tt.fields.UnregisteredJobs,
			}
			got := c.GetStatusJSON()
			assert.NotNil(t, got)
		})
	}
}
