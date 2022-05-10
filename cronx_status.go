package cronx

import (
	"time"

	"github.com/robfig/cron/v3"
)

//go:generate gomodifytags -all --skip-unexported -w -file cronx_status.go -remove-tags db,json
//go:generate gomodifytags -all --skip-unexported -w -file cronx_status.go -add-tags db,json

// StatusData defines current job status.
type StatusData struct {
	// ID is unique per job.
	ID cron.EntryID `db:"id" json:"id"`
	// Job defines current job.
	Job *Job `db:"job" json:"job"`
	// Next defines the next schedule to execute current job.
	Next time.Time `db:"next" json:"next"`
	// Prev defines the last run of the current job.
	Prev time.Time `db:"prev" json:"prev"`
}

type StatusPageData struct {
	Data []StatusData `db:"data" json:"data"`
	Sort Sort         `db:"sort" json:"sort"`
}

type Sort struct {
	Query   string            `db:"query" json:"query"`
	Columns map[string]string `db:"columns" json:"columns"`
}
