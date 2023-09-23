-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "transaction" (
    id              UUID            PRIMARY KEY,
    description     VARCHAR(50)     NOT NULL,
    created_at      TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    amount          NUMERIC         NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "transaction";
-- +goose StatementEnd
