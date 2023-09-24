-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "uuid" CHAR(26) NOT NULL UNIQUE,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "username" VARCHAR(255) NOT NULL UNIQUE,
  "encrypted_password" BYTEA NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "users";
-- +goose StatementEnd
