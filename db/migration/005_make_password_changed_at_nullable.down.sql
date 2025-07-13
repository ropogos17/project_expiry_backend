-- Migration 005 Rollback: Make password_changed_at field NOT NULL again in users table
-- This migration adds back the NOT NULL constraint to password_changed_at column
-- Note: This rollback will fail if there are NULL values in the column

-- Update any NULL values to current timestamp before adding NOT NULL constraint
UPDATE "users" SET "password_changed_at" = NOW() WHERE "password_changed_at" IS NULL;

-- Add back the NOT NULL constraint
ALTER TABLE "users" ALTER COLUMN "password_changed_at" SET NOT NULL;