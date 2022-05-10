package cronx

import (
	"net/url"

	"github.com/rizalgowandy/cronx/pagination"
	"github.com/rizalgowandy/cronx/storage"
)

//go:generate gomodifytags -all --skip-unexported -w -file cronx_history.go -remove-tags db,json
//go:generate gomodifytags -all --skip-unexported -w -file cronx_history.go -add-tags db,json

type HistoryData struct {
	Data       []storage.History   `db:"data"       json:"data"`
	Pagination pagination.Response `db:"pagination" json:"pagination"`
}

func generateURI(param map[string]string) *string {
	val := url.Values{}
	for k, v := range param {
		val.Add(k, v)
	}
	res := val.Encode()
	return &res
}
