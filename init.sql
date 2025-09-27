CREATE TABLE IF NOT EXISTS "users" (
  "id" integer PRIMARY KEY,
  "username" varchar,
  "password_hash" varchar,
  "email" varchar,
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "posts" (
  "id" integer PRIMARY KEY,
  "user_id" integer,
  "title" varchar,
  "content" varchar,
  "created_at" timestamp,
  "updated_at" timestamp,
  "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "refresh_tokens" (
  "id" integer PRIMARY KEY,
  "user_id" integer,
  "token" text UNIQUE,
  "expires_at" TIMESTAMP
);

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "refresh_tokens" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

