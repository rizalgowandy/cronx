package cronx

import (
	"github.com/rizalgowandy/cronx/storage"
	"github.com/rizalgowandy/gdk/pkg/pagination"
)

//go:generate gomodifytags -all --quiet -w -file cronx_history.go -clear-tags
//go:generate gomodifytags -all --quiet --skip-unexported -w -file cronx_history.go -add-tags json

type HistoryPageData struct {
	Data       []storage.History `json:"data"`
	Pagination Response          `json:"pagination"`
	Sort       pagination.Sort   `json:"sort"`
}
