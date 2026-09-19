\connect jobs_service

ALTER TABLE jobs
ADD COLUMN offer_url TEXT NOT NULL DEFAULT '';
