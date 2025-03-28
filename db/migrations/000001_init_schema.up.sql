CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "email" varchar UNIQUE NOT NULL,
  "full_name" varchar NOT NULL,
  "hashed_password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "payments" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "amount" bigint NOT NULL,
  "status" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "payments" ("owner");

ALTER TABLE "payments" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");