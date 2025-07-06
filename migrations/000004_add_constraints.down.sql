-- Drop triggers
DROP TRIGGER IF EXISTS update_shopping_list_items_updated_at ON shopping_list_items;
DROP TRIGGER IF EXISTS update_items_updated_at ON items;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop constraints
ALTER TABLE storage DROP CONSTRAINT IF EXISTS check_shelves_positive;
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_email_format;
ALTER TABLE users DROP CONSTRAINT IF EXISTS check_login_attempts;
ALTER TABLE shares DROP CONSTRAINT IF EXISTS unique_share_per_user_storage;