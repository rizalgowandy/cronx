package cronx

import (
	"context"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
)

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
	}
	return
}

type JobMetadata struct {
	EntryID    cron.EntryID `json:"entry_id"`
	Wave       int64        `json:"wave"`
	TotalWave  int64        `json:"total_wave"`
	IsLastWave bool         `json:"is_last_wave"`
}

type Job struct {
	JobMetadata

	Name    string     `json:"name"`
	Status  StatusCode `json:"status"`
	Latency string     `json:"latency"`
	Error   string     `json:"error"`

	manager *Manager
	inner   JobItf
	status  uint32
	running sync.Mutex
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
	ctx := context.Background()

	// Lock current process.
	j.running.Lock()
	defer j.running.Unlock()

	// Set job metadata.
	ctx = SetJobMetadata(ctx, j.JobMetadata)

	// Update job status as running.
	atomic.StoreUint32(&j.status, statusRunning)
	j.UpdateStatus()

	// Run the job.
	if err := j.manager.Interceptor(ctx, j, func(ctx context.Context, job *Job) error {
		return job.inner.Run(ctx)
	}); err != nil {
		j.Error = err.Error()
		atomic.StoreUint32(&j.status, statusError)
	} else {
		atomic.StoreUint32(&j.status, statusIdle)
	}

	// Record time needed to execute the whole process.
	j.Latency = time.Since(start).String()

	// Update job status after running.
	j.UpdateStatus()
}
