-- Migration 006: Add deleted_at, created_at, and updated_at columns to tables that don't have them
-- This migration adds standard timestamp columns for better data tracking

-- Add timestamp columns to roles table
ALTER TABLE "roles" ADD COLUMN "deleted_at" TIMESTAMP WITH TIME ZONE;
ALTER TABLE "roles" ADD COLUMN "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();
ALTER TABLE "roles" ADD COLUMN "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();

-- Add timestamp columns to user_statuses table
ALTER TABLE "user_statuses" ADD COLUMN "deleted_at" TIMESTAMP WITH TIME ZONE;
ALTER TABLE "user_statuses" ADD COLUMN "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();
ALTER TABLE "user_statuses" ADD COLUMN "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();

-- Add timestamp columns to categories table
ALTER TABLE "categories" ADD COLUMN "deleted_at" TIMESTAMP WITH TIME ZONE;
ALTER TABLE "categories" ADD COLUMN "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();
ALTER TABLE "categories" ADD COLUMN "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();

-- Add timestamp columns to storages table
ALTER TABLE "storages" ADD COLUMN "deleted_at" TIMESTAMP WITH TIME ZONE;
ALTER TABLE "storages" ADD COLUMN "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();
ALTER TABLE "storages" ADD COLUMN "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();

-- Add timestamp columns to announcements table
ALTER TABLE "announcements" ADD COLUMN "deleted_at" TIMESTAMP WITH TIME ZONE;
ALTER TABLE "announcements" ADD COLUMN "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();
ALTER TABLE "announcements" ADD COLUMN "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();

-- Add timestamp columns to languages table
ALTER TABLE "languages" ADD COLUMN "deleted_at" TIMESTAMP WITH TIME ZONE;
ALTER TABLE "languages" ADD COLUMN "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();
ALTER TABLE "languages" ADD COLUMN "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();

-- Add timestamp columns to action_types table
ALTER TABLE "action_types" ADD COLUMN "deleted_at" TIMESTAMP WITH TIME ZONE;
ALTER TABLE "action_types" ADD COLUMN "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();
ALTER TABLE "action_types" ADD COLUMN "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW();

-- Add deleted_at column to refresh_tokens table (it already has timestamp columns but different names)
ALTER TABLE "refresh_tokens" ADD COLUMN "deleted_at" TIMESTAMP WITH TIME ZONE;