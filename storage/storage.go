package storage

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/jsonx"
	"github.com/rizalgowandy/gdk/pkg/sortx"
)

//go:generate gomodifytags -all --quiet -w -file storage.go -clear-tags
//go:generate gomodifytags -all --quiet --skip-unexported -w -file storage.go -add-tags db,json
//go:generate gomodifytags --quiet --skip-unexported -w -file storage.go -struct HistoryMetadata -add-tags json -add-options json=omitempty
//go:generate gomodifytags --quiet --skip-unexported -w -file storage.go -struct ErrorDetail -add-tags json -add-options json=omitempty

type Client interface {
	WriteHistory(ctx context.Context, req *History) error
	ReadHistories(ctx context.Context, req *HistoryFilter) ([]History, error)
}

type History struct {
	ID          int64           `db:"id"           json:"id"`
	CreatedAt   time.Time       `db:"created_at"   json:"created_at"`
	Name        string          `db:"name"         json:"name"`
	Status      string          `db:"status"       json:"status"`
	StatusCode  int64           `db:"status_code"  json:"status_code"`
	StartedAt   time.Time       `db:"started_at"   json:"started_at"`
	FinishedAt  time.Time       `db:"finished_at"  json:"finished_at"`
	Latency     int64           `db:"latency"      json:"latency"`
	LatencyText string          `db:"latency_text" json:"latency_text"`
	Error       ErrorDetail     `db:"error"        json:"error"`
	Metadata    HistoryMetadata `db:"metadata"     json:"metadata"`
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
	Sorts         sortx.Sorts `db:"sorts"          json:"sorts"`
	Limit         int         `db:"limit"          json:"limit"`
	StartingAfter *int64      `db:"starting_after" json:"starting_after"`
	EndingBefore  *int64      `db:"ending_before"  json:"ending_before"`
}
