package storage

import (
	"context"
)

// NewNoopClient returns a no operation client.
func NewNoopClient() *NoopClient {
	return &NoopClient{}
}

type NoopClient struct{}

func (n NoopClient) WriteHistory(_ context.Context, _ *History) error {
	return nil
}

func (n NoopClient) ReadHistories(_ context.Context, _ *HistoryFilter) ([]History, error) {
	return nil, nil
}
