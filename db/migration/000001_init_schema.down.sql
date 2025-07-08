-- Drop foreign keys first to avoid dependency errors
ALTER TABLE "Log" DROP CONSTRAINT IF EXISTS "Log_user_id_fkey";
ALTER TABLE "Log" DROP CONSTRAINT IF EXISTS "Log_action_type_fkey";
ALTER TABLE "Share" DROP CONSTRAINT IF EXISTS "Share_target_user_id_fkey";
ALTER TABLE "Share" DROP CONSTRAINT IF EXISTS "Share_storage_id_fkey";
ALTER TABLE "Share" DROP CONSTRAINT IF EXISTS "Share_source_user_id_fkey";
ALTER TABLE "ShoppingListItem" DROP CONSTRAINT IF EXISTS "ShoppingListItem_last_modified_fkey";
ALTER TABLE "Storage" DROP CONSTRAINT IF EXISTS "Storage_owner_fkey";
ALTER TABLE "Item" DROP CONSTRAINT IF EXISTS "Item_last_modified_fkey";
ALTER TABLE "Item" DROP CONSTRAINT IF EXISTS "Item_storage_id_fkey";
ALTER TABLE "Item" DROP CONSTRAINT IF EXISTS "Item_category_id_fkey";
ALTER TABLE "User" DROP CONSTRAINT IF EXISTS "User_status_id_fkey";
ALTER TABLE "User" DROP CONSTRAINT IF EXISTS "User_role_id_fkey";

-- Drop tables (in reverse order of creation to respect dependencies)
DROP TABLE IF EXISTS "ActionType";
DROP TABLE IF EXISTS "Log";
DROP TABLE IF EXISTS "Languages";
DROP TABLE IF EXISTS "Share";
DROP TABLE IF EXISTS "ShoppingListItem";
DROP TABLE IF EXISTS "Announcement";
DROP TABLE IF EXISTS "Storage";
DROP TABLE IF EXISTS "Category";
DROP TABLE IF EXISTS "Item";
DROP TABLE IF EXISTS "UserStatus";
DROP TABLE IF EXISTS "Role";
DROP TABLE IF EXISTS "User";
