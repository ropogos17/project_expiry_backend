-- Drop foreign key constraints (optional if using CASCADE in DROP TABLE)
ALTER TABLE "logs" DROP CONSTRAINT IF EXISTS "logs_user_id_fkey";
ALTER TABLE "logs" DROP CONSTRAINT IF EXISTS "logs_action_type_fkey";

ALTER TABLE "shares" DROP CONSTRAINT IF EXISTS "shares_target_user_id_fkey";
ALTER TABLE "shares" DROP CONSTRAINT IF EXISTS "shares_storage_id_fkey";
ALTER TABLE "shares" DROP CONSTRAINT IF EXISTS "shares_source_user_id_fkey";

ALTER TABLE "shopping_list_items" DROP CONSTRAINT IF EXISTS "shopping_list_items_last_modified_fkey";

ALTER TABLE "storages" DROP CONSTRAINT IF EXISTS "storages_owner_fkey";

ALTER TABLE "items" DROP CONSTRAINT IF EXISTS "items_last_modified_fkey";
ALTER TABLE "items" DROP CONSTRAINT IF EXISTS "items_storage_id_fkey";
ALTER TABLE "items" DROP CONSTRAINT IF EXISTS "items_category_id_fkey";

ALTER TABLE "users" DROP CONSTRAINT IF EXISTS "users_status_id_fkey";
ALTER TABLE "users" DROP CONSTRAINT IF EXISTS "users_role_id_fkey";

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS "logs";
DROP TABLE IF EXISTS "action_types";
DROP TABLE IF EXISTS "languages";
DROP TABLE IF EXISTS "shares";
DROP TABLE IF EXISTS "shopping_list_items";
DROP TABLE IF EXISTS "announcements";
DROP TABLE IF EXISTS "storages";
DROP TABLE IF EXISTS "categories";
DROP TABLE IF EXISTS "items";
DROP TABLE IF EXISTS "user_statuses";
DROP TABLE IF EXISTS "roles";
DROP TABLE IF EXISTS "users";
