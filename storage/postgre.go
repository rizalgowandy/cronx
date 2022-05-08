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
			machine_id,
			created_at,
			entry_id,
			name,
			started_at,
			finished_at,
			latency
		)
		VALUES (
		   $1,
		   $2,
		   $3,
		   $4,
		   $5,
		   $6,
		   $7,
		   $8
		)
		;
	`

	_, err = pool.Exec(
		ctx,
		query,
		param.ID,
		param.MachineID,
		param.CreatedAt,
		param.EntryID,
		param.Name,
		param.StartedAt,
		param.FinishedAt,
		param.Latency,
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

func (p *PostgreClient) ReadHistories(ctx context.Context, param pagination.Request) (ReadHistoriesRes, error) {
	panic("implement me")
}
