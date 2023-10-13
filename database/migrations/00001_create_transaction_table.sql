-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
    id              TEXT            PRIMARY KEY,
    description     VARCHAR(50)     NOT NULL,
    date            DATE            NOT NULL,
    amount          NUMERIC         NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions;
-- +goose StatementEnd
