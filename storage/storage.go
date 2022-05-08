package storage

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/rizalgowandy/gdk/pkg/pagination"
)

type Client interface {
	WriteHistory(ctx context.Context, param *History) error
	ReadHistories(ctx context.Context, param pagination.Request) (ReadHistoriesRes, error)
}

type History struct {
	ID         string    `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	Name       string    `db:"name"`
	Status     string    `db:"status"`
	StatusCode int64     `db:"status_code"`
	StartedAt  time.Time `db:"started_at"`
	FinishedAt time.Time `db:"finished_at"`
	Latency    int64     `db:"latency"`
	Metadata   HistoryMetadata
}

type HistoryMetadata struct {
	MachineID  string `json:"machine_id"`
	EntryID    int64  `json:"entry_id"`
	Wave       int64  `json:"wave"`
	TotalWave  int64  `json:"total_wave"`
	IsLastWave bool   `json:"is_last_wave"`
}

func (h HistoryMetadata) Value() (driver.Value, error) {
	return json.Marshal(h)
}

func (h *HistoryMetadata) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &h)
}

type ReadHistoriesRes struct {
	Data       []History
	Pagination pagination.Response
}
