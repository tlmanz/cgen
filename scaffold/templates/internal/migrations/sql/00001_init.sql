-- +goose Up
CREATE TABLE accounts (
    id       SERIAL PRIMARY KEY,
    owner    VARCHAR(255)   NOT NULL,
    currency VARCHAR(3)     NOT NULL,
    balance  NUMERIC(15, 2) NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS accounts;
