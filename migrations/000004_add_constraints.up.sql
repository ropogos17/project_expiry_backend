-- Add unique constraints
ALTER TABLE shares ADD CONSTRAINT unique_share_per_user_storage
    UNIQUE (source_user_id, target_user_id, storage_id);

-- Add check constraints
ALTER TABLE users ADD CONSTRAINT check_login_attempts
    CHECK (login_attempts >= 0);

ALTER TABLE users ADD CONSTRAINT check_email_format
    CHECK (email ~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

ALTER TABLE storage ADD CONSTRAINT check_shelves_positive
    CHECK (shelves IS NULL OR shelves > 0);

-- Add triggers for updated_at columns
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_items_updated_at
    BEFORE UPDATE ON items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_shopping_list_items_updated_at
    BEFORE UPDATE ON shopping_list_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();