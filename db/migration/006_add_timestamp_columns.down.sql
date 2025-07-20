-- Migration 006 Rollback: Remove timestamp columns from tables
-- This migration removes the added timestamp columns

-- Remove timestamp columns from roles table
ALTER TABLE "roles" DROP COLUMN IF EXISTS "updated_at";
ALTER TABLE "roles" DROP COLUMN IF EXISTS "created_at";
ALTER TABLE "roles" DROP COLUMN IF EXISTS "deleted_at";

-- Remove timestamp columns from user_statuses table
ALTER TABLE "user_statuses" DROP COLUMN IF EXISTS "updated_at";
ALTER TABLE "user_statuses" DROP COLUMN IF EXISTS "created_at";
ALTER TABLE "user_statuses" DROP COLUMN IF EXISTS "deleted_at";

-- Remove timestamp columns from categories table
ALTER TABLE "categories" DROP COLUMN IF EXISTS "updated_at";
ALTER TABLE "categories" DROP COLUMN IF EXISTS "created_at";
ALTER TABLE "categories" DROP COLUMN IF EXISTS "deleted_at";

-- Remove timestamp columns from storages table
ALTER TABLE "storages" DROP COLUMN IF EXISTS "updated_at";
ALTER TABLE "storages" DROP COLUMN IF EXISTS "created_at";
ALTER TABLE "storages" DROP COLUMN IF EXISTS "deleted_at";

-- Remove timestamp columns from announcements table
ALTER TABLE "announcements" DROP COLUMN IF EXISTS "updated_at";
ALTER TABLE "announcements" DROP COLUMN IF EXISTS "created_at";
ALTER TABLE "announcements" DROP COLUMN IF EXISTS "deleted_at";

-- Remove timestamp columns from languages table
ALTER TABLE "languages" DROP COLUMN IF EXISTS "updated_at";
ALTER TABLE "languages" DROP COLUMN IF EXISTS "created_at";
ALTER TABLE "languages" DROP COLUMN IF EXISTS "deleted_at";

-- Remove timestamp columns from action_types table
ALTER TABLE "action_types" DROP COLUMN IF EXISTS "updated_at";
ALTER TABLE "action_types" DROP COLUMN IF EXISTS "created_at";
ALTER TABLE "action_types" DROP COLUMN IF EXISTS "deleted_at";

-- Remove deleted_at column from refresh_tokens table
ALTER TABLE "refresh_tokens" DROP COLUMN IF EXISTS "deleted_at";