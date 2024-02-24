-- +goose Up
CREATE TABLE s3 (
    id UUID PRIMARY KEY,
    bucket TEXT UNIQUE NOT NULL,
    tags TEXT NOT NULL,
    bucket_domain_name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE s3;