-- +goose Up
-- modify "businesses" table
ALTER TABLE "businesses" ADD CONSTRAINT "chk_businesses_revenue" CHECK (revenue >= 0), ADD COLUMN "revenue" bigint NOT NULL DEFAULT 0;
-- modify "services" table
ALTER TABLE "services" ADD CONSTRAINT "chk_services_charge" CHECK (charge > 0), ADD COLUMN "charge" integer NOT NULL;

-- +goose Down
-- reverse: modify "services" table
ALTER TABLE "services" DROP COLUMN "charge", DROP CONSTRAINT "chk_services_charge";
-- reverse: modify "businesses" table
ALTER TABLE "businesses" DROP COLUMN "revenue", DROP CONSTRAINT "chk_businesses_revenue";
