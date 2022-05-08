CREATE TABLE cronx_histories (
	id          TEXT            NOT NULL
		CONSTRAINT cronx_histories_id_pk
			PRIMARY KEY,
	machine_id  TEXT DEFAULT '' NOT NULL,
	created_at  TIMESTAMPTZ     NOT NULL,
	entry_id    INT8            NOT NULL,
	name        TEXT DEFAULT '' NOT NULL,
	started_at  TIMESTAMPTZ     NOT NULL,
	finished_at TIMESTAMPTZ     NOT NULL,
	latency     INT8            NOT NULL
);

CREATE INDEX cronx_histories_created_at_index
	ON cronx_histories(created_at DESC);
