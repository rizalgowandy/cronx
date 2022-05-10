CREATE TABLE cronx_histories (
	id           TEXT               NOT NULL
		CONSTRAINT cronx_histories_id_pk
			PRIMARY KEY,
	created_at   TIMESTAMPTZ        NOT NULL,
	name         TEXT               NOT NULL,
	status       TEXT               NOT NULL,
	status_code  INT8               NOT NULL,
	started_at   TIMESTAMPTZ        NOT NULL,
	finished_at  TIMESTAMPTZ        NOT NULL,
	latency      INT8               NOT NULL,
	latency_text TEXT               NOT NULL,
	error        JSONB DEFAULT '{}' NOT NULL,
	metadata     JSONB DEFAULT '{}' NOT NULL
);

CREATE INDEX cronx_histories_created_at_index
	ON cronx_histories(created_at DESC);
