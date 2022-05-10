package cronx

import (
	"net/url"

	"github.com/rizalgowandy/cronx/storage"
	"github.com/rizalgowandy/gdk/pkg/pagination"
)

//go:generate gomodifytags -all --skip-unexported -w -file cronx_history.go -remove-tags db,json
//go:generate gomodifytags -all --skip-unexported -w -file cronx_history.go -add-tags db,json -add-options json=omitempty

type HistoryData struct {
	Data       []storage.History   `db:"data"       json:"data,omitempty"`
	Pagination pagination.Response `db:"pagination" json:"pagination,omitempty"`
}

func generateURI(param map[string]string) *string {
	val := url.Values{}
	for k, v := range param {
		val.Add(k, v)
	}
	res := val.Encode()
	return &res
}
