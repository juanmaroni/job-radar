\connect jobs_service

CREATE TABLE jobs (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    source TEXT NOT NULL,
    source_id TEXT NOT NULL,
    company TEXT NOT NULL,
    offer_url TEXT NOT NULL DEFAULT '',
    role TEXT NOT NULL,
    location TEXT,
    salary_min INTEGER,
    salary_max INTEGER,
    description TEXT,
    score INTEGER NOT NULL,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(source, source_id)
);

GRANT ALL PRIVILEGES ON TABLE jobs TO jobs_service;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO jobs_service;
