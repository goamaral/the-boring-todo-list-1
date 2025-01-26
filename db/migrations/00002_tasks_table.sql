-- +goose Up
-- +goose StatementBegin
CREATE TABLE "tasks" (
  "id" SERIAL PRIMARY KEY,
  "uuid" UUID NOT NULL UNIQUE,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "title" VARCHAR(255) NOT NULL,
  "done_at" TIMESTAMP,
  "author_id" SERIAL REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "tasks";
-- +goose StatementEnd
