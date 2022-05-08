CREATE TABLE cronx_jobs (
	name       TEXT        NOT NULL
		CONSTRAINT cronx_jobs_name_pk
			PRIMARY KEY,
	created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX cronx_jobs_created_at_index
	ON cronx_jobs(created_at DESC);
