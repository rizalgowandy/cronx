package storage

import (
	"context"
	"sync"

	"github.com/rizalgowandy/gdk/pkg/errorx/v2"
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

func (p *PostgreClient) WriteHistory(ctx context.Context, param *History) error {
	fields := errorx.Fields{tags.Fields: param}

	pool, err := p.db.Get(ctx)
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
		param.ID,
		param.CreatedAt,
		param.Name,
		param.Status,
		param.StatusCode,
		param.StartedAt,
		param.FinishedAt,
		param.Latency,
		param.Error,
		param.Metadata,
	)
	if err != nil {
		return errorx.E(err, fields)
	}

	_, exists := p.jobs.Load(param.Name)
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
			param.Name,
			param.CreatedAt,
		)
		if err != nil {
			return errorx.E(err, fields)
		}
		p.jobs.Store(param.Name, true)
	}

	return nil
}

func (p *PostgreClient) ReadHistories(
	ctx context.Context,
	param pagination.Request,
) (ReadHistoriesRes, error) {
	panic("implement me")
}
