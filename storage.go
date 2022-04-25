package cronx

import "context"

type Storage interface {
	WriteHistory(ctx context.Context, param WriteHistoryParam) error
	ReadHistory(ctx context.Context, filter ReadHistoryFilter) (ReadHistoryRes, error)
}

type WriteHistoryParam struct {
}

type ReadHistoryFilter struct {
}

type ReadHistoryRes struct {
}
