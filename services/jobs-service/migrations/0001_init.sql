CREATE TABLE jobs (
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    source     TEXT NOT NULL,
    source_id  TEXT NOT NULL,
    company    TEXT NOT NULL,
    role       TEXT NOT NULL,
    score      INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(source, source_id)
);
