package storage

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/jsonx"
)

//go:generate gomodifytags -all --skip-unexported -w -file storage.go -remove-tags db,json
//go:generate gomodifytags -all --skip-unexported -w -file storage.go -add-tags db,json -add-options json=omitempty

type Client interface {
	WriteHistory(ctx context.Context, req *History) error
	ReadHistories(ctx context.Context, req *HistoryFilter) ([]History, error)
}

type History struct {
	ID         string          `db:"id"          json:"id,omitempty"`
	CreatedAt  time.Time       `db:"created_at"  json:"created_at,omitempty"`
	Name       string          `db:"name"        json:"name,omitempty"`
	Status     string          `db:"status"      json:"status,omitempty"`
	StatusCode int64           `db:"status_code" json:"status_code,omitempty"`
	StartedAt  time.Time       `db:"started_at"  json:"started_at,omitempty"`
	FinishedAt time.Time       `db:"finished_at" json:"finished_at,omitempty"`
	Latency    int64           `db:"latency"     json:"latency,omitempty"`
	Error      ErrorDetail     `db:"error"       json:"error,omitempty"`
	Metadata   HistoryMetadata `db:"metadata"    json:"metadata,omitempty"`
}

type HistoryMetadata struct {
	MachineID  string `db:"machine_id"   json:"machine_id,omitempty"`
	EntryID    int64  `db:"entry_id"     json:"entry_id,omitempty"`
	Wave       int64  `db:"wave"         json:"wave,omitempty"`
	TotalWave  int64  `db:"total_wave"   json:"total_wave,omitempty"`
	IsLastWave bool   `db:"is_last_wave" json:"is_last_wave,omitempty"`
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
	Err          string              `db:"err"           json:"err,omitempty"`
	Code         errorx.Code         `db:"code"          json:"code,omitempty"`
	Fields       errorx.Fields       `db:"fields"        json:"fields,omitempty"`
	OpTraces     []errorx.Op         `db:"op_traces"     json:"op_traces,omitempty"`
	Message      errorx.Message      `db:"message"       json:"message,omitempty"`
	Line         errorx.Line         `db:"line"          json:"line,omitempty"`
	MetricStatus errorx.MetricStatus `db:"metric_status" json:"metric_status,omitempty"`
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

type HistoryFilter struct {
	Order         string  `db:"order"          json:"order,omitempty"`
	Limit         int     `db:"limit"          json:"limit,omitempty"`
	StartingAfter *string `db:"starting_after" json:"starting_after,omitempty"`
	EndingBefore  *string `db:"ending_before"  json:"ending_before,omitempty"`
}
