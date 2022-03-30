package cronx

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
)

func TestNewCommandController(t *testing.T) {
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
			got := NewCommandController(tt.args.config, tt.args.interceptors)
			assert.NotNil(t, got)
		})
	}
}

func TestCommandController_StatusData(t *testing.T) {
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
			name:   "Commander is nil",
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
			c := &CommandController{
				Commander:        tt.fields.Commander,
				Interceptor:      tt.fields.Interceptor,
				Parser:           tt.fields.Parser,
				UnregisteredJobs: tt.fields.UnregisteredJobs,
			}
			if got := c.StatusData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StatusData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandController_StatusJSON(t *testing.T) {
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
			name:   "Commander is nil",
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
			c := &CommandController{
				Commander:        tt.fields.Commander,
				Interceptor:      tt.fields.Interceptor,
				Parser:           tt.fields.Parser,
				UnregisteredJobs: tt.fields.UnregisteredJobs,
			}
			got := c.StatusJSON()
			assert.NotNil(t, got)
		})
	}
}

func TestCommandController_Info(t *testing.T) {
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
			c := &CommandController{
				Commander:        tt.fields.Commander,
				Interceptor:      tt.fields.Interceptor,
				Parser:           tt.fields.Parser,
				UnregisteredJobs: tt.fields.UnregisteredJobs,
				Address:          tt.fields.Address,
				Location:         tt.fields.Location,
				CreatedTime:      tt.fields.CreatedTime,
			}
			got := c.Info()
			assert.NotNil(t, got)
		})
	}
}
