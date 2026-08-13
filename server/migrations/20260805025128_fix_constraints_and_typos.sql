-- +goose Up
-- modify "businesses" table
ALTER TABLE "businesses" ALTER COLUMN "status" SET DEFAULT 'OPEN';

-- +goose Down
-- reverse: modify "businesses" table
ALTER TABLE "businesses" ALTER COLUMN "status" DROP DEFAULT;
