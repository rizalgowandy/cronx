package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
	"github.com/rizalgowandy/gdk/pkg/sortx"
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
		sq = sq.Where("id > ?", *req.StartingAfter)
	}
	if req.EndingBefore != nil {
		sq = sq.FromSelect(
			sq.From("cronx_histories").
				Where("id < ?", *req.EndingBefore).
				OrderBy(createOrderBy(req.Sorts, true)),
			"before",
		)
	}
	sq = sq.OrderBy(createOrderBy(req.Sorts, false))

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

func createOrderBy(sorts []sortx.Sort, reverse bool) string {
	var result string
	for _, v := range sorts {
		cur := v.Order
		if reverse {
			switch v.Order {
			case sortx.OrderAscending:
				cur = sortx.OrderDescending
			case sortx.OrderDescending:
				cur = sortx.OrderAscending
			}
		}
		if result != "" {
			result += ", "
		}
		result += fmt.Sprintf("%s %s", v.Key, cur)
	}
	return result
}
