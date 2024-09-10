-- Reverse the foreign key addition on "accounts" referencing "users"
ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS accounts_owner_fkey;

-- Drop the unique index on "accounts" ("owner", "currency")
DROP INDEX IF EXISTS accounts_owner_currency_idx;

-- Optionally drop the accounts table (only if you want to fully remove it)
-- DROP TABLE IF EXISTS "accounts";

-- Drop the users table
DROP TABLE IF EXISTS "users";
