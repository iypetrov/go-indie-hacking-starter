-- +goose NO TRANSACTION
-- +goose Up
PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;
PRAGMA busy_timeout = 5000;

-- +goose Down
PRAGMA journal_mode = DELETE;
PRAGMA foreign_keys = OFF;
PRAGMA busy_timeout = 0;
