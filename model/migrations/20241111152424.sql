-- Create "quotes" table
CREATE TABLE "public"."quotes" (
  "id" bigserial NOT NULL,
  "customer_id" bigint NOT NULL,
  "description" text NOT NULL,
  "status" text NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "unique_customer_id_description" to table: "quotes"
CREATE UNIQUE INDEX "unique_customer_id_description" ON "public"."quotes" ("customer_id", "description");
-- Create "items" table
CREATE TABLE "public"."items" (
  "quote_id" bigint NOT NULL,
  "item_id" text NOT NULL,
  "segment" text NOT NULL,
  "price" numeric NOT NULL,
  "tax" numeric NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  PRIMARY KEY ("quote_id", "item_id", "segment"),
  CONSTRAINT "fk_quotes_items" FOREIGN KEY ("quote_id") REFERENCES "public"."quotes" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
