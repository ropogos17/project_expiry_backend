-- Remove initial data
DELETE FROM action_types WHERE name IN (
                                        'login', 'logout', 'create_item', 'update_item', 'delete_item',
                                        'create_storage', 'update_storage', 'delete_storage',
                                        'share_storage', 'unshare_storage', 'create_shopping_item',
                                        'update_shopping_item', 'delete_shopping_item', 'user_created',
                                        'user_updated', 'user_deleted', 'password_changed',
                                        'email_verified', 'phone_verified', 'two_factor_enabled',
                                        'two_factor_disabled'
    );

DELETE FROM categories WHERE name IN (
                                      'Food', 'Beverages', 'Household', 'Personal Care',
                                      'Medicine', 'Electronics', 'Tools', 'Other'
    );

DELETE FROM status WHERE name IN ('active', 'inactive', 'suspended', 'pending');
DELETE FROM roles WHERE name IN ('admin', 'user', 'moderator');