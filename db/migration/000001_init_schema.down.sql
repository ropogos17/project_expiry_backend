-- ============================================================
-- 000001_init_schema.down.sql
-- Drop in reverse dependency order
-- ============================================================

DROP INDEX IF EXISTS idx_logs_created_at;
DROP INDEX IF EXISTS idx_logs_action_type_id;
DROP INDEX IF EXISTS idx_logs_user_id;

DROP INDEX IF EXISTS idx_announcements_deleted_at;
DROP INDEX IF EXISTS idx_announcements_created_by;

DROP INDEX IF EXISTS idx_shares_storage_id;
DROP INDEX IF EXISTS idx_shares_target_user_id;
DROP INDEX IF EXISTS idx_shares_source_user_id;

DROP INDEX IF EXISTS idx_items_deleted_at;
DROP INDEX IF EXISTS idx_items_expires_at;
DROP INDEX IF EXISTS idx_items_category_id;
DROP INDEX IF EXISTS idx_items_storage_id;

DROP INDEX IF EXISTS idx_storages_deleted_at;
DROP INDEX IF EXISTS idx_storages_owner;

DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;

DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_status_id;
DROP INDEX IF EXISTS idx_users_role_id;

DROP TABLE IF EXISTS logs;
DROP TABLE IF EXISTS announcements;
DROP TABLE IF EXISTS shares;
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS storages;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS action_types;
DROP TABLE IF EXISTS user_statuses;
DROP TABLE IF EXISTS roles;