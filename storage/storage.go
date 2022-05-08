package storage

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/jsonx"
	"github.com/rizalgowandy/gdk/pkg/pagination"
)

type Client interface {
	WriteHistory(ctx context.Context, param *History) error
	ReadHistories(ctx context.Context, param pagination.Request) (ReadHistoriesRes, error)
}

type History struct {
	ID         string          `db:"id"`
	CreatedAt  time.Time       `db:"created_at"`
	Name       string          `db:"name"`
	Status     string          `db:"status"`
	StatusCode int64           `db:"status_code"`
	StartedAt  time.Time       `db:"started_at"`
	FinishedAt time.Time       `db:"finished_at"`
	Latency    int64           `db:"latency"`
	Error      ErrorDetail     `db:"error"`
	Metadata   HistoryMetadata `db:"metadata"`
}

type HistoryMetadata struct {
	MachineID  string `json:"machine_id,omitempty"`
	EntryID    int64  `json:"entry_id,omitempty"`
	Wave       int64  `json:"wave,omitempty"`
	TotalWave  int64  `json:"total_wave,omitempty"`
	IsLastWave bool   `json:"is_last_wave,omitempty"`
}

func (h *HistoryMetadata) Value() (driver.Value, error) {
	return jsonx.Marshal(h)
}

func (h *HistoryMetadata) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errorx.E("type assertion to []byte failed")
	}

	return jsonx.Unmarshal(b, &h)
}

type ErrorDetail struct {
	Err          string              `json:"err,omitempty"`
	Code         errorx.Code         `json:"code,omitempty"`
	Fields       errorx.Fields       `json:"fields,omitempty"`
	OpTraces     []errorx.Op         `json:"op,omitempty"`
	Message      errorx.Message      `json:"message,omitempty"`
	Line         errorx.Line         `json:"line,omitempty"`
	MetricStatus errorx.MetricStatus `json:"metric_status,omitempty"`
}

func (e *ErrorDetail) Value() (driver.Value, error) {
	return jsonx.Marshal(e)
}

func (e *ErrorDetail) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errorx.E("type assertion to []byte failed")
	}

	return jsonx.Unmarshal(b, &e)
}

type ReadHistoriesRes struct {
	Data       []History
	Pagination pagination.Response
}
