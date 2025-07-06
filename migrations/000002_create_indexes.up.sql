-- Indexes for users table
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role_id ON users(role_id);
CREATE INDEX idx_users_status_id ON users(status_id);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

-- Indexes for items table
CREATE INDEX idx_items_category_id ON items(category_id);
CREATE INDEX idx_items_storage_id ON items(storage_id);
CREATE INDEX idx_items_expires_at ON items(expires_at);
CREATE INDEX idx_items_stored_at ON items(stored_at);
CREATE INDEX idx_items_is_deleted ON items(is_deleted);
CREATE INDEX idx_items_last_modified ON items(last_modified);

-- Indexes for storage table
CREATE INDEX idx_storage_owner ON storage(owner);

-- Indexes for shares table
CREATE INDEX idx_shares_source_user_id ON shares(source_user_id);
CREATE INDEX idx_shares_target_user_id ON shares(target_user_id);
CREATE INDEX idx_shares_storage_id ON shares(storage_id);
CREATE INDEX idx_shares_shared_at ON shares(shared_at);

-- Indexes for shopping_list_items table
CREATE INDEX idx_shopping_list_items_last_modified ON shopping_list_items(last_modified);
CREATE INDEX idx_shopping_list_items_created_at ON shopping_list_items(created_at);
CREATE INDEX idx_shopping_list_items_is_deleted ON shopping_list_items(is_deleted);

-- Indexes for logs table
CREATE INDEX idx_logs_user_id ON logs(user_id);
CREATE INDEX idx_logs_action_type ON logs(action_type);
CREATE INDEX idx_logs_created_at ON logs(created_at);