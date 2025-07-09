-- Migration 002 Rollback: Rename action_type_id column back to action_type in logs table
-- This migration reverts the column name from action_type_id back to action_type

-- Drop the foreign key constraint
ALTER TABLE "logs" DROP CONSTRAINT IF EXISTS "logs_action_type_id_fkey";

-- Rename the column back to original name
ALTER TABLE "logs" RENAME COLUMN "action_type_id" TO "action_type";

-- Re-add the foreign key constraint with the original column name
ALTER TABLE "logs" ADD CONSTRAINT "logs_action_type_fkey" FOREIGN KEY ("action_type") REFERENCES "action_types" ("id");