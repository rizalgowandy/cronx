package cronx

import "context"

type contextKey string

// Context key for standardized context value.
const (
	// CtxKeyJobMetadata is context for cron job metadata.
	CtxKeyJobMetadata = contextKey("cron-job-metadata")
)

// GetJobMetadata returns job metadata from current context, and status if it exists or not.
func GetJobMetadata(ctx context.Context) (JobMetadata, bool) {
	if ctx == nil {
		return JobMetadata{}, false
	}
	val := ctx.Value(CtxKeyJobMetadata)

	meta, ok := (val).(JobMetadata)
	if !ok {
		return JobMetadata{}, false
	}

	return meta, true
}

// SetJobMetadata stores current job metadata inside current context.
func SetJobMetadata(ctx context.Context, meta JobMetadata) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, CtxKeyJobMetadata, meta)
}
