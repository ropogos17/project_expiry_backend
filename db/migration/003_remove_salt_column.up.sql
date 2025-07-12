-- Migration 003: Remove salt column from users table
-- This migration removes the salt column as it's no longer needed

ALTER TABLE "users" DROP COLUMN IF EXISTS "salt";