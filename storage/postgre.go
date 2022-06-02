package storage

import (
	"context"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
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
			created_at,
			name,
			status,
			status_code,
			started_at,
			finished_at,
			latency,
			latency_text,
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
		req.CreatedAt,
		req.Name,
		req.Status,
		req.StatusCode,
		req.StartedAt,
		req.FinishedAt,
		req.Latency,
		req.LatencyText,
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

func (p *PostgreClient) ReadHistories(ctx context.Context, req *HistoryFilter) ([]History, error) {
	fields := errorx.Fields{tags.Request: req}

	pool, err := p.db.GetReader(ctx)
	if err != nil {
		return nil, errorx.E(err, fields)
	}

	sq := squirrel.
		Select(
			"id",
			"created_at",
			"name",
			"status",
			"status_code",
			"started_at",
			"finished_at",
			"latency",
			"latency_text",
			"error",
			"metadata",
		).
		From("cronx_histories").
		Limit(uint64(req.Limit)).
		PlaceholderFormat(squirrel.Dollar)

	if req.StartingAfter != nil {
		if !req.Sorts.Desc() {
			sq = sq.Where("id > ?", *req.StartingAfter)
		} else {
			sq = sq.Where("id < ?", *req.StartingAfter)
		}
	}
	if req.EndingBefore != nil {
		before := sq.OrderBy(req.Sorts.OrderBy(true))

		if !req.Sorts.Desc() {
			before = before.Where("id < ?", *req.EndingBefore)
		} else {
			before = before.Where("id > ?", *req.EndingBefore)
		}

		sq = sq.FromSelect(before, "before")
	}
	sq = sq.OrderBy(req.Sorts.OrderBy())

	query, args, err := sq.ToSql()
	if err != nil {
		return nil, errorx.E(err, fields)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		if errorx.Match(err, pgx.ErrNoRows) {
			return nil, errorx.E(err, fields, errorx.CodeNotFound)
		}
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
			&cur.LatencyText,
			&cur.Error,
			&cur.Metadata,
		); err != nil {
			return nil, errorx.E(err, fields)
		}

		data = append(data, cur)
	}

	return data, nil
}
