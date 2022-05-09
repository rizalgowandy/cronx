package storage

import (
	"context"

	"github.com/rizalgowandy/gdk/pkg/pagination"
)

// NewNoopClient returns a no operation client.
func NewNoopClient() *NoopClient {
	return &NoopClient{}
}

type NoopClient struct{}

func (n NoopClient) WriteHistory(_ context.Context, _ *History) error {
	return nil
}

func (n NoopClient) ReadHistories(
	_ context.Context,
	_ pagination.Request,
) (*ReadHistoriesRes, error) {
	return &ReadHistoriesRes{}, nil
}
