package cronx

import "context"

// Func is a type to allow callers to wrap a raw func.
// Example:
//
//	manager.Schedule("@every 5m", cronx.Func(myFunc))
type Func func(ctx context.Context) error

func (f Func) Run(ctx context.Context) error {
	return f(ctx)
}

// NewFuncJob creates a wrapper type for a raw func with name.
func NewFuncJob(name string, cmd Func) *FuncJob {
	return &FuncJob{
		name: name,
		cmd:  cmd,
	}
}

// FuncJob is a type to allow callers to wrap a raw func with name.
// Example:
//
//	manager.ScheduleFunc("@every 5m", "random name", myFunc)
type FuncJob struct {
	name string
	cmd  Func
}

func (f FuncJob) Run(ctx context.Context) error {
	return f.cmd(ctx)
}
