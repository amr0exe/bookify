-- +goose Up
-- create "services" table
CREATE TABLE "services" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "business_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "desc" character varying(255) NOT NULL,
  "duration" integer NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_businesses_service" FOREIGN KEY ("business_id") REFERENCES "businesses" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_services_duration" CHECK (duration > 0)
);
-- create index "idx_services_business_id" to table: "services"
CREATE INDEX "idx_services_business_id" ON "services" ("business_id");
-- create index "idx_services_deleted_at" to table: "services"
CREATE INDEX "idx_services_deleted_at" ON "services" ("deleted_at");

-- +goose Down
-- reverse: create index "idx_services_deleted_at" to table: "services"
DROP INDEX "idx_services_deleted_at";
-- reverse: create index "idx_services_business_id" to table: "services"
DROP INDEX "idx_services_business_id";
-- reverse: create "services" table
DROP TABLE "services";
