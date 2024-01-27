-- Drop the foreign key constraint and unique index on "accounts"
ALTER TABLE "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";
DROP INDEX IF EXISTS "accounts_owner_idx";

-- Drop the "users" table
DROP TABLE IF EXISTS "users";
