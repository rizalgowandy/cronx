package cronx

import (
	"net/url"
	"strconv"

	"github.com/rizalgowandy/gdk/pkg/converter"
	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
)

//go:generate gomodifytags -all --quiet -w -file pagination.go -clear-tags
//go:generate gomodifytags -all --quiet --skip-unexported -w -file pagination.go -add-tags query,form,json,xml

// Request is a parameter to return list of data with pagination.
// Request is optional, most fields automatically filled by system.
// If you already have a response with pagination,
// you can generate pagination request directly to traverse next or prev page.
type Request struct {
	url url.URL

	// Sort of the resources in the response e.g. sort=id:desc,created_at:desc
	// Sort is optional.
	Sort string `query:"sort"           form:"sort"           json:"sort"           xml:"sort"`
	// Limit number of results per call. Accepted values: 1 - 100. Default 10
	// Limit is optional.
	Limit int `query:"limit"          form:"limit"          json:"limit"          xml:"limit"`
	// StartingAfter is a cursor for use in pagination.
	// StartingAfter is a resource ID that defines your place in the list.
	// StartingAfter is optional.
	StartingAfter *int64 `query:"starting_after" form:"starting_after" json:"starting_after" xml:"starting_after"`
	// EndingBefore is cursor for use in pagination.
	// EndingBefore is a resource ID that defines your place in the list.
	// EndingBefore is optional.
	EndingBefore *int64 `query:"ending_before"  form:"ending_before"  json:"ending_before"  xml:"ending_before"`
}

func (r *Request) Validate() error {
	if r.url == (url.URL{}) {
		return errorx.E("url cannot be empty")
	}
	if r.Sort == "" {
		r.Sort = "id"
	}
	if r.Limit == 0 {
		r.Limit = 10
	}
	return nil
}

func (r *Request) QueryParams() map[string]string {
	res := map[string]string{}
	if r.Sort != "" {
		res["sort"] = r.Sort
	}
	if r.Limit > 0 {
		res["limit"] = strconv.Itoa(r.Limit)
	}
	if r.StartingAfter != nil {
		res["starting_after"] = converter.String(*r.StartingAfter)
	}
	if r.EndingBefore != nil {
		res["ending_before"] = converter.String(*r.EndingBefore)
	}
	return res
}

func (r *Request) URI(req *url.URL) *string {
	val := &url.Values{}
	for k, v := range r.QueryParams() {
		val.Add(k, v)
	}
	req.RawQuery = val.Encode()
	uri := req.RequestURI()
	return &uri
}

type Response struct {
	Sort          string  `query:"sort"           form:"sort"           json:"sort"           xml:"sort"`
	StartingAfter *int64  `query:"starting_after" form:"starting_after" json:"starting_after" xml:"starting_after"`
	EndingBefore  *int64  `query:"ending_before"  form:"ending_before"  json:"ending_before"  xml:"ending_before"`
	Total         int     `query:"total"          form:"total"          json:"total"          xml:"total"`
	Yielded       int     `query:"yielded"        form:"yielded"        json:"yielded"        xml:"yielded"`
	Limit         int     `query:"limit"          form:"limit"          json:"limit"          xml:"limit"`
	PreviousURI   *string `query:"previous_uri"   form:"previous_uri"   json:"previous_uri"   xml:"previous_uri"`
	NextURI       *string `query:"next_uri"       form:"next_uri"       json:"next_uri"       xml:"next_uri"`
	// CursorRange returns cursors for starting after and ending before.
	// Format: [starting_after, ending_before].
	CursorRange []int64 `query:"cursor_range"   form:"cursor_range"   json:"cursor_range"   xml:"cursor_range"`
}

// HasPrevPage returns true if prev page exists and can be traversed.
func (r *Response) HasPrevPage() bool {
	return r.PreviousURI != nil
}

// HasNextPage returns true if next page exists and can be traversed.
func (r *Response) HasNextPage() bool {
	return r.NextURI != nil
}

// PrevPageCursor returns cursor to be used as ending before value.
func (r *Response) PrevPageCursor() *int64 {
	if len(r.CursorRange) < 1 {
		return nil
	}
	return &r.CursorRange[0]
}

// NextPageCursor returns cursor to be used as starting after value.
func (r *Response) NextPageCursor() *int64 {
	if len(r.CursorRange) < 2 {
		return nil
	}
	return &r.CursorRange[1]
}

// PrevPageRequest returns pagination request for the prev page result.
func (r *Response) PrevPageRequest() *Request {
	return &Request{
		Sort:          r.Sort,
		Limit:         r.Limit,
		StartingAfter: nil,
		EndingBefore:  r.PrevPageCursor(),
	}
}

// NextPageRequest returns pagination request for the next page result.
func (r *Response) NextPageRequest() *Request {
	return &Request{
		Sort:          r.Sort,
		Limit:         r.Limit,
		StartingAfter: r.NextPageCursor(),
		EndingBefore:  nil,
	}
}

type Sort struct {
	Query   string            `query:"query"   form:"query"   json:"query"   xml:"query"`
	Columns map[string]string `query:"columns" form:"columns" json:"columns" xml:"columns"`
}
