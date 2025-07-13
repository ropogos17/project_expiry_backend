-- Migration 005: Make password_changed_at field nullable in users table
-- This migration removes the NOT NULL constraint from password_changed_at column

ALTER TABLE "users" ALTER COLUMN "password_changed_at" DROP NOT NULL;