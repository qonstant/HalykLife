CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

CREATE UNIQUE INDEX ON "accounts" ("phone_number");

-- ALTER TABLE "accounts" ADD CONSTRAINT "accounts_owner_idx" UNIQUE ("phone_number");
