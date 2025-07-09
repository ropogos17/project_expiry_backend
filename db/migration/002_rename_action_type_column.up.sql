-- Migration 002: Rename action_type column to action_type_id in logs table
-- This migration renames the column from action_type to action_type_id for better naming consistency

-- Drop the existing foreign key constraint
ALTER TABLE "logs" DROP CONSTRAINT IF EXISTS "logs_action_type_fkey";

-- Rename the column
ALTER TABLE "logs" RENAME COLUMN "action_type" TO "action_type_id";

-- Re-add the foreign key constraint with the new column name
ALTER TABLE "logs" ADD CONSTRAINT "logs_action_type_id_fkey" FOREIGN KEY ("action_type_id") REFERENCES "action_types" ("id");