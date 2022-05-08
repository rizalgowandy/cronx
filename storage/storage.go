package storage

import "context"

type Client interface {
	WriteHistory(ctx context.Context, param WriteHistoryParam) error
	ReadHistories(ctx context.Context, filter ReadHistoriesFilter) (ReadHistoriesRes, error)
}

type WriteHistoryParam struct {
}

type ReadHistoriesFilter struct {
	Limit int64
}

type ReadHistoriesRes struct {
}
