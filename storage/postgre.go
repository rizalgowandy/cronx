package storage

import (
	"context"
	"net/url"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/logx"
	"github.com/rizalgowandy/gdk/pkg/pagination"
	"github.com/rizalgowandy/gdk/pkg/storage/database"
	"github.com/rizalgowandy/gdk/pkg/tags"
)

func NewPostgreClient(db database.PostgreClientItf) *PostgreClient {
	return &PostgreClient{
		db:   db,
		jobs: sync.Map{},
	}
}

type PostgreClient struct {
	db   database.PostgreClientItf
	jobs sync.Map
}

func (p *PostgreClient) WriteHistory(ctx context.Context, req *History) error {
	fields := errorx.Fields{tags.Request: req}

	pool, err := p.db.GetWriter(ctx)
	if err != nil {
		return errorx.E(err, fields)
	}

	query := `
		INSERT INTO cronx_histories (
			id,
			created_at,
			name,
			status,
			status_code,
			started_at,
			finished_at,
			latency,
			error,
			metadata
		)
		VALUES (
		   $1,
		   $2,
		   $3,
		   $4,
		   $5,
		   $6,
		   $7,
		   $8,
		   $9,
		   $10
		)
		;
	`

	_, err = pool.Exec(
		ctx,
		query,
		req.ID,
		req.CreatedAt,
		req.Name,
		req.Status,
		req.StatusCode,
		req.StartedAt,
		req.FinishedAt,
		req.Latency,
		req.Error,
		req.Metadata,
	)
	if err != nil {
		return errorx.E(err, fields)
	}

	_, exists := p.jobs.Load(req.Name)
	if !exists {
		query = `
			INSERT INTO cronx_jobs (
				name,
				created_at
			)
			VALUES (
			   $1,
			   $2
			)
			ON CONFLICT DO NOTHING
			;
		`

		_, err = pool.Exec(
			ctx,
			query,
			req.Name,
			req.CreatedAt,
		)
		if err != nil {
			return errorx.E(err, fields)
		}
		p.jobs.Store(req.Name, true)
	}

	return nil
}

func (p *PostgreClient) ReadHistories(
	ctx context.Context,
	req pagination.Request,
) (*ReadHistoriesRes, error) {
	fields := errorx.Fields{tags.Request: req}

	pool, err := p.db.GetReader(ctx)
	if err != nil {
		return nil, errorx.E(err, fields)
	}

	var sq squirrel.SelectBuilder
	if req.Order == "" {
		req.Order = "created_at DESC"
	}
	if req.Limit == 0 {
		req.Limit = 25
	}

	sq = squirrel.
		Select(
			"id",
			"created_at",
			"name",
			"status",
			"status_code",
			"started_at",
			"finished_at",
			"latency",
			"error",
			"metadata",
		).
		From("cronx_histories").
		OrderBy(req.Order).
		Limit(uint64(req.Limit)).
		PlaceholderFormat(squirrel.Dollar)

	if req.StartingAfter != nil {
		sq = sq.Where("id > ?", *req.StartingAfter)
	}
	if req.EndingBefore != nil {
		sq = sq.Where("id < ?", *req.EndingBefore)
	}

	query, args, err := sq.ToSql()
	if err != nil {
		return nil, errorx.E(err, fields)
	}

	logx.DBG(ctx, args, query)

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, errorx.E(err, fields)
	}
	defer rows.Close()

	var data []History
	for rows.Next() {
		var cur History

		if err := rows.Scan(
			&cur.ID,
			&cur.CreatedAt,
			&cur.Name,
			&cur.Status,
			&cur.StatusCode,
			&cur.StartedAt,
			&cur.FinishedAt,
			&cur.Latency,
			&cur.Error,
			&cur.Metadata,
		); err != nil {
			return nil, errorx.E(err, fields)
		}

		data = append(data, cur)
	}

	page := pagination.Response{
		Order:         req.Order,
		StartingAfter: req.StartingAfter,
		EndingBefore:  req.EndingBefore,
		Total:         len(data),
		Yielded:       len(data),
		Limit:         req.Limit,
		PreviousURI:   nil,
		NextURI:       nil,
		CursorRange:   nil,
	}
	if len(data) > 0 {
		page.CursorRange = []string{
			data[0].ID,
			data[len(data)-1].ID,
		}
		page.NextURI = generateURI(page.NextPageRequest().QueryParams())
		if req.StartingAfter != nil {
			page.PreviousURI = generateURI(page.PrevPageRequest().QueryParams())
		}
	}

	return &ReadHistoriesRes{
		Data:       data,
		Pagination: page,
	}, nil
}

func generateURI(param map[string]string) *string {
	val := url.Values{}
	for k, v := range param {
		val.Add(k, v)
	}
	res := val.Encode()
	return &res
}
