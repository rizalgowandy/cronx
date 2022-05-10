package cronx

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rizalgowandy/cronx/storage"
	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/logx"
	"github.com/rizalgowandy/gdk/pkg/netx"
	"github.com/robfig/cron/v3"
)

//go:generate gomodifytags -all --skip-unexported -w -file job.go -remove-tags db,json
//go:generate gomodifytags -all --skip-unexported -w -file job.go -add-tags db,json

type JobItf interface {
	Run(ctx context.Context) error
}

// NewJob creates a new job with default status and name.
func NewJob(manager *Manager, job JobItf, waveNumber, totalWave int64) *Job {
	return &Job{
		manager: manager,
		JobMetadata: JobMetadata{
			EntryID:    0,
			Wave:       waveNumber,
			TotalWave:  totalWave,
			IsLastWave: waveNumber == totalWave,
		},
		Name:    GetJobName(job),
		Status:  StatusCodeUp,
		Latency: "",
		Error:   "",
		inner:   job,
		status:  statusUp,
		running: sync.Mutex{},
	}
}

// GetJobName return the Job name by reflect the job
func GetJobName(job JobItf) (name string) {
	name = reflect.TypeOf(job).Name()
	if name == "" {
		name = reflect.TypeOf(job).Elem().Name()
	}
	if name == "" {
		name = reflect.TypeOf(job).String()
	}
	if name == "Func" {
		name = "(nameless)"
		return
	}
	if name == "FuncJob" {
		fj, ok := job.(*FuncJob)
		if !ok {
			name = "(nameless)"
			return
		}
		name = fj.name
	}
	return
}

type JobMetadata struct {
	EntryID    cron.EntryID `db:"entry_id" json:"entry_id"`
	Wave       int64        `db:"wave" json:"wave"`
	TotalWave  int64        `db:"total_wave" json:"total_wave"`
	IsLastWave bool         `db:"is_last_wave" json:"is_last_wave"`
}

type Job struct {
	JobMetadata

	Name    string     `db:"name" json:"name"`
	Status  StatusCode `db:"status" json:"status"`
	Latency string     `db:"latency" json:"latency"`
	Error   string     `db:"error" json:"error"`

	manager *Manager
	inner   JobItf
	status  uint32
	running sync.Mutex
	latency int64
	err     error
}

// UpdateStatus updates the current job status to the latest.
func (j *Job) UpdateStatus() StatusCode {
	switch atomic.LoadUint32(&j.status) {
	case statusRunning:
		j.Status = StatusCodeRunning
	case statusIdle:
		j.Status = StatusCodeIdle
	case statusDown:
		j.Status = StatusCodeDown
	case statusError:
		j.Status = StatusCodeError
	default:
		j.Status = StatusCodeUp
	}
	return j.Status
}

// Run executes the current job operation.
func (j *Job) Run() {
	start := time.Now()
	ctx := logx.NewContext()

	// Lock current process.
	j.running.Lock()
	defer j.running.Unlock()

	// Set job metadata.
	ctx = SetJobMetadata(ctx, j.JobMetadata)

	// Update job status as running.
	atomic.StoreUint32(&j.status, statusRunning)
	j.UpdateStatus()

	// Run the job.
	if err := j.manager.interceptor(ctx, j, func(ctx context.Context, job *Job) error {
		return job.inner.Run(ctx)
	}); err != nil {
		j.err = err
		j.Error = err.Error()
		atomic.StoreUint32(&j.status, statusError)
	} else {
		atomic.StoreUint32(&j.status, statusIdle)
	}

	// Record time needed to execute the whole process.
	finish := time.Now()
	latency := time.Since(start)
	j.latency = latency.Nanoseconds()
	j.Latency = latency.String()

	// Update job status after running.
	j.UpdateStatus()

	// Record history.
	j.RecordHistory(ctx, start, finish)
}

func (j *Job) RecordHistory(ctx context.Context, start, finish time.Time) {
	history := &storage.History{
		ID:         logx.GetRequestID(ctx),
		CreatedAt:  time.Now(),
		Name:       j.Name,
		Status:     j.Status.String(),
		StatusCode: int64(j.status),
		StartedAt:  start,
		FinishedAt: finish,
		Latency:    j.latency,
		Metadata: storage.HistoryMetadata{
			MachineID: netx.GetIPv4(),
			EntryID:   int64(j.JobMetadata.EntryID),
		},
	}

	// Only add wave information for job with multiple wave.
	if j.JobMetadata.TotalWave > 1 {
		history.Metadata.Wave = j.JobMetadata.Wave
		history.Metadata.TotalWave = j.JobMetadata.TotalWave
		history.Metadata.IsLastWave = j.JobMetadata.IsLastWave
	}

	// Add error detail.
	if j.err != nil {
		history.Error.Err = j.err.Error()
		if e, ok := j.err.(*errorx.Error); ok {
			history.Error = storage.ErrorDetail{
				Err:          e.Err.Error(),
				Code:         e.Code,
				Fields:       e.Fields,
				OpTraces:     e.OpTraces,
				Message:      e.Message,
				Line:         e.Line,
				MetricStatus: e.MetricStatus,
			}
		}
	}

	if err := j.manager.storage.WriteHistory(ctx, history); err != nil {
		logx.ERR(ctx, errorx.E(err), "write history must success")
	}
}
