-- Migration 004: Add refresh_tokens table and update timestamp columns to include WITH TIME ZONE
-- This migration adds the refresh_tokens table and modifies existing timestamp columns

-- Create refresh_tokens table
CREATE TABLE "refresh_tokens" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id" INTEGER NOT NULL REFERENCES "users"("id") ON DELETE CASCADE,
    "issued_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "expires_at" TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Update timestamp columns in users table
ALTER TABLE "users" ALTER COLUMN "last_login" TYPE TIMESTAMP WITH TIME ZONE;
ALTER TABLE "users" ALTER COLUMN "locked_until" TYPE TIMESTAMP WITH TIME ZONE;
ALTER TABLE "users" ALTER COLUMN "password_changed_at" TYPE TIMESTAMP WITH TIME ZONE;
ALTER TABLE "users" ALTER COLUMN "deleted_at" TYPE TIMESTAMP WITH TIME ZONE;
ALTER TABLE "users" ALTER COLUMN "created_at" TYPE TIMESTAMP WITH TIME ZONE;
ALTER TABLE "users" ALTER COLUMN "updated_at" TYPE TIMESTAMP WITH TIME ZONE;

-- Update timestamp columns in items table
ALTER TABLE "items" ALTER COLUMN "stored_at" TYPE TIMESTAMP WITH TIME ZONE;
ALTER TABLE "items" ALTER COLUMN "updated_at" TYPE TIMESTAMP WITH TIME ZONE;
ALTER TABLE "items" ALTER COLUMN "expires_at" TYPE TIMESTAMP WITH TIME ZONE;

-- Update timestamp columns in shopping_list_items table
ALTER TABLE "shopping_list_items" ALTER COLUMN "created_at" TYPE TIMESTAMP WITH TIME ZONE;
ALTER TABLE "shopping_list_items" ALTER COLUMN "updated_at" TYPE TIMESTAMP WITH TIME ZONE;

-- Update timestamp columns in shares table
ALTER TABLE "shares" ALTER COLUMN "shared_at" TYPE TIMESTAMP WITH TIME ZONE;

-- Update timestamp columns in logs table
ALTER TABLE "logs" ALTER COLUMN "created_at" TYPE TIMESTAMP WITH TIME ZONE;