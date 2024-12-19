-- +goose Up
CREATE TABLE IF NOT EXISTS mailing_list (
    id BLOB PRIMARY KEY, 
    email TEXT NOT NULL UNIQUE,
    last_sent_at TEXT,
    created_at TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS mailing_list;
