-- +goose Up
CREATE TABLE IF NOT EXISTS mailing_list (
    id INTEGER PRIMARY KEY, 
    email TEXT NOT NULL UNIQUE,
    last_sent_at TEXT
);

-- +goose Down
DROP TABLE IF EXISTS mailing_list;
