package cronx

// StatusCode describes current job status.
type StatusCode string

func (s StatusCode) String() string {
	return string(s)
}

const (
	// StatusCodeUp describes that current job has just been created.
	StatusCodeUp StatusCode = "UP"
	// StatusCodeIdle describes that current job is waiting for next execution time.
	StatusCodeIdle StatusCode = "IDLE"
	// StatusCodeRunning describes that current job is currently running.
	StatusCodeRunning StatusCode = "RUNNING"
	// StatusCodeDown describes that current job has failed to be registered.
	StatusCodeDown StatusCode = "DOWN"
	// StatusCodeError describes that last run has failed.
	StatusCodeError StatusCode = "ERROR"

	statusDown    uint32 = 0
	statusUp      uint32 = 1
	statusIdle    uint32 = 2
	statusRunning uint32 = 3
	statusError   uint32 = 4
)
