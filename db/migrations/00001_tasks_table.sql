-- +goose Up
-- +goose StatementBegin
CREATE TABLE "tasks" (
  "id" CHAR(26) PRIMARY KEY,
  "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "title" VARCHAR(255) NOT NULL,
  "completed_at" TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "tasks";
-- +goose StatementEnd
