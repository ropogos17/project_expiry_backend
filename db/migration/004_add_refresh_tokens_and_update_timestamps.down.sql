-- Migration 004 Rollback: Remove refresh_tokens table and revert timestamp columns
-- This migration removes the refresh_tokens table and reverts timestamp columns back to timestamp without time zone

-- Drop refresh_tokens table
DROP TABLE IF EXISTS "refresh_tokens";

-- Revert timestamp columns in users table
ALTER TABLE "users" ALTER COLUMN "last_login" TYPE TIMESTAMP;
ALTER TABLE "users" ALTER COLUMN "locked_until" TYPE TIMESTAMP;
ALTER TABLE "users" ALTER COLUMN "password_changed_at" TYPE TIMESTAMP;
ALTER TABLE "users" ALTER COLUMN "deleted_at" TYPE TIMESTAMP;
ALTER TABLE "users" ALTER COLUMN "created_at" TYPE TIMESTAMP;
ALTER TABLE "users" ALTER COLUMN "updated_at" TYPE TIMESTAMP;

-- Revert timestamp columns in items table
ALTER TABLE "items" ALTER COLUMN "stored_at" TYPE TIMESTAMP;
ALTER TABLE "items" ALTER COLUMN "updated_at" TYPE TIMESTAMP;
ALTER TABLE "items" ALTER COLUMN "expires_at" TYPE TIMESTAMP;

-- Revert timestamp columns in shopping_list_items table
ALTER TABLE "shopping_list_items" ALTER COLUMN "created_at" TYPE TIMESTAMP;
ALTER TABLE "shopping_list_items" ALTER COLUMN "updated_at" TYPE TIMESTAMP;

-- Revert timestamp columns in shares table
ALTER TABLE "shares" ALTER COLUMN "shared_at" TYPE TIMESTAMP;

-- Revert timestamp columns in logs table
ALTER TABLE "logs" ALTER COLUMN "created_at" TYPE TIMESTAMP;