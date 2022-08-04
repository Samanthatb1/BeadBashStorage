CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "full_name" varchar NOT NULL,
  "total_orders" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "orders" (
  "order_id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "username" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "purchase_amount" float NOT NULL,
  "purchased_item" varchar NOT NULL,
  "shipping_location" varchar NOT NULL,
  "currency" varchar NOT NULL,
  "date_ordered" varchar NOT NULL
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "orders" ("account_id");

CREATE INDEX ON "orders" ("username");

CREATE INDEX ON "orders" ("shipping_location");

COMMENT ON COLUMN "orders"."purchase_amount" IS 'must be positive';

ALTER TABLE "orders" ADD FOREIGN KEY ("account_id") REFERENCES "users" ("id");

ALTER TABLE "orders" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
