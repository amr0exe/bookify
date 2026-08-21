-- +goose Up
-- create "appointments" table
CREATE TABLE "appointments" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "consumer_id" uuid NOT NULL,
  "business_id" uuid NOT NULL,
  "service_id" uuid NOT NULL,
  "status" character varying(20) NOT NULL DEFAULT 'created',
  "duration" integer NOT NULL,
  "remarks" text NULL,
  "scheduled_at" timestamptz NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_appointments_duration" CHECK (duration > 0)
);
-- create index "idx_appointments_business_id" to table: "appointments"
CREATE INDEX "idx_appointments_business_id" ON "appointments" ("business_id");
-- create index "idx_appointments_consumer_id" to table: "appointments"
CREATE INDEX "idx_appointments_consumer_id" ON "appointments" ("consumer_id");
-- create index "idx_appointments_deleted_at" to table: "appointments"
CREATE INDEX "idx_appointments_deleted_at" ON "appointments" ("deleted_at");
-- create index "idx_appointments_service_id" to table: "appointments"
CREATE INDEX "idx_appointments_service_id" ON "appointments" ("service_id");

-- +goose Down
-- reverse: create index "idx_appointments_service_id" to table: "appointments"
DROP INDEX "idx_appointments_service_id";
-- reverse: create index "idx_appointments_deleted_at" to table: "appointments"
DROP INDEX "idx_appointments_deleted_at";
-- reverse: create index "idx_appointments_consumer_id" to table: "appointments"
DROP INDEX "idx_appointments_consumer_id";
-- reverse: create index "idx_appointments_business_id" to table: "appointments"
DROP INDEX "idx_appointments_business_id";
-- reverse: create "appointments" table
DROP TABLE "appointments";
