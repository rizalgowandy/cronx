package cronx

import (
	"time"

	"github.com/robfig/cron/v3"
)

//go:generate gomodifytags -all --skip-unexported -w -file cronx_status.go -remove-tags db,json
//go:generate gomodifytags -all --skip-unexported -w -file cronx_status.go -add-tags db,json -add-options json=omitempty

// StatusData defines current job status.
type StatusData struct {
	// ID is unique per job.
	ID cron.EntryID `db:"id" json:"id,omitempty"`
	// Job defines current job.
	Job *Job `db:"job" json:"job,omitempty"`
	// Next defines the next schedule to execute current job.
	Next time.Time `db:"next" json:"next,omitempty"`
	// Prev defines the last run of the current job.
	Prev time.Time `db:"prev" json:"prev,omitempty"`
}

type StatusPageData struct {
	Data []StatusData `db:"data" json:"data,omitempty"`
	Sort Sort         `db:"sort" json:"sort,omitempty"`
}

type Sort struct {
	Query   string            `db:"query" json:"query,omitempty"`
	Columns map[string]string `db:"columns" json:"columns,omitempty"`
}
