-- Create the users table with a unique constraint on full_name
CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar UNIQUE NOT NULL,  -- Make full_name unique
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamp NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamp NOT NULL DEFAULT (now())
);

-- Add a foreign key on accounts referencing the full_name
ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("full_name");

-- Create a unique index on the combination of owner and currency
CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
