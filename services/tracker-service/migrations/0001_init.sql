CREATE TYPE status_type AS ENUM ('applied', 'interview', 'offer', 'rejected');

CREATE TABLE applications (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    company VARCHAR(64) NOT NULL,
    role VARCHAR(64) NOT NULL,
    offer_url VARCHAR(255) NOT NULL,
    status status_type NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
