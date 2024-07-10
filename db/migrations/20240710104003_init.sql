-- +goose Up
CREATE TABLE "user" (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL
);

CREATE TABLE booking (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES "user"(id),
  startDate TIMESTAMP NOT NULL,
  endDate TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS booking;
DROP TABLE IF EXISTS "user";
