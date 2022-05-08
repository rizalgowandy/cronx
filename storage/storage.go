package storage

import (
	"context"
	"time"

	"github.com/rizalgowandy/gdk/pkg/pagination"
)

type Client interface {
	WriteHistory(ctx context.Context, param *History) error
	ReadHistories(ctx context.Context, param pagination.Request) (ReadHistoriesRes, error)
}

type History struct {
	ID         string    `db:"id"`
	MachineID  string    `db:"machine_id"`
	CreatedAt  time.Time `db:"created_at"`
	EntryID    int64     `db:"entry_id"`
	Name       string    `db:"name"`
	StartedAt  time.Time `db:"started_at"`
	FinishedAt time.Time `db:"finished_at"`
	Latency    int64     `db:"latency"`
}

type ReadHistoriesRes struct {
	Data       []History
	Pagination pagination.Response
}
