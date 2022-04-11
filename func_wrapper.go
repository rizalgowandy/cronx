package cronx

import "context"

// Func is a type to allow callers to wrap a raw func.
// Example:
//	manager.Schedule("@every 5m", cronx.Func(myFunc))
type Func func(ctx context.Context) error

func (r Func) Run(ctx context.Context) error {
	return r(ctx)
}
