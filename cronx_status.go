package cronx

import (
	"time"

	"github.com/robfig/cron/v3"
)

//go:generate gomodifytags -all --quiet -w -file cronx_status.go -clear-tags
//go:generate gomodifytags -all --quiet --skip-unexported -w -file cronx_status.go -add-tags json

// StatusData defines current job status.
type StatusData struct {
	// ID is unique per job.
	ID cron.EntryID `json:"id"`
	// Job defines current job.
	Job *Job `json:"job"`
	// Next defines the next schedule to execute current job.
	Next time.Time `json:"next"`
	// Prev defines the last run of the current job.
	Prev time.Time `json:"prev"`
}

type StatusPageData struct {
	Data []StatusData `json:"data"`
	Sort Sort         `json:"sort"`
}
