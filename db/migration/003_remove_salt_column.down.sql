-- Migration 003 Rollback: Add salt column back to users table
-- This migration adds the salt column back to the users table

ALTER TABLE "users" ADD COLUMN "salt" varchar NOT NULL DEFAULT '';