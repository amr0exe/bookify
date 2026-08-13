-- +goose Up
-- create "businesses" table
CREATE TABLE "businesses" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "account_id" uuid NOT NULL,
  "name" character varying(254) NOT NULL,
  "description" text NULL,
  "phone" character varying(50) NOT NULL,
  "address" text NOT NULL,
  "latitude" numeric(10,8) NULL,
  "longitude" numeric(11,8) NULL,
  "status" character varying(20) NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_accounts_business" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_businesses_account_id" to table: "businesses"
CREATE UNIQUE INDEX "idx_businesses_account_id" ON "businesses" ("account_id");
-- create index "idx_businesses_deleted_at" to table: "businesses"
CREATE INDEX "idx_businesses_deleted_at" ON "businesses" ("deleted_at");
-- create "consumers" table
CREATE TABLE "consumers" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "account_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "display_name" character varying(255) NULL,
  "phone" character varying(60) NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_accounts_consumer" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create index "idx_consumers_account_id" to table: "consumers"
CREATE UNIQUE INDEX "idx_consumers_account_id" ON "consumers" ("account_id");
-- create index "idx_consumers_deleted_at" to table: "consumers"
CREATE INDEX "idx_consumers_deleted_at" ON "consumers" ("deleted_at");

-- +goose Down
-- reverse: create index "idx_consumers_deleted_at" to table: "consumers"
DROP INDEX "idx_consumers_deleted_at";
-- reverse: create index "idx_consumers_account_id" to table: "consumers"
DROP INDEX "idx_consumers_account_id";
-- reverse: create "consumers" table
DROP TABLE "consumers";
-- reverse: create index "idx_businesses_deleted_at" to table: "businesses"
DROP INDEX "idx_businesses_deleted_at";
-- reverse: create index "idx_businesses_account_id" to table: "businesses"
DROP INDEX "idx_businesses_account_id";
-- reverse: create "businesses" table
DROP TABLE "businesses";
