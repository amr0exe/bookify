-- +goose Up
-- create "refresh_tokens" table
CREATE TABLE "refresh_tokens" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "account_id" uuid NOT NULL,
  "token_hash" text NOT NULL,
  "is_revoked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_refresh_tokens_account_id" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_account_id" ON "refresh_tokens" ("account_id");
-- create index "idx_refresh_tokens_token_hash" to table: "refresh_tokens"
CREATE UNIQUE INDEX "idx_refresh_tokens_token_hash" ON "refresh_tokens" ("token_hash");

-- +goose Down
-- reverse: create index "idx_refresh_tokens_token_hash" to table: "refresh_tokens"
DROP INDEX "idx_refresh_tokens_token_hash";
-- reverse: create index "idx_refresh_tokens_account_id" to table: "refresh_tokens"
DROP INDEX "idx_refresh_tokens_account_id";
-- reverse: create "refresh_tokens" table
DROP TABLE "refresh_tokens";
