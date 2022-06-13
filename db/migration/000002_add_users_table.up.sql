CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "email_changed_at" timestamp NOT NULL DEFAULT ('0001-01-01 00:00:00'),
  "created_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username") ON DELETE CASCADE;
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");
-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency")