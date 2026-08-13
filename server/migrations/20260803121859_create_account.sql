-- +goose Up
-- create "accounts" table
CREATE TABLE "public"."accounts" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "email" character varying(255) NOT NULL,
  "pass_hash" character varying(255) NOT NULL,
  "role" character varying(20) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_accounts_role" CHECK ((role)::text = ANY ((ARRAY['BUSINESS'::character varying, 'CONSUMER'::character varying])::text[]))
);
-- create index "idx_accounts_deleted_at" to table: "accounts"
CREATE INDEX "idx_accounts_deleted_at" ON "public"."accounts" ("deleted_at");
-- create index "idx_accounts_email" to table: "accounts"
CREATE UNIQUE INDEX "idx_accounts_email" ON "public"."accounts" ("email");

-- +goose Down
-- reverse: create index "idx_accounts_email" to table: "accounts"
DROP INDEX "public"."idx_accounts_email";
-- reverse: create index "idx_accounts_deleted_at" to table: "accounts"
DROP INDEX "public"."idx_accounts_deleted_at";
-- reverse: create "accounts" table
DROP TABLE "public"."accounts";
