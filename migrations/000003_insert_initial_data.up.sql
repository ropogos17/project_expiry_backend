-- Insert default roles
INSERT INTO roles (name) VALUES
                             ('admin'),
                             ('user'),
                             ('moderator');

-- Insert default status
INSERT INTO status (name) VALUES
                              ('active'),
                              ('inactive'),
                              ('suspended'),
                              ('pending');

-- Insert default categories
INSERT INTO categories (name) VALUES
                                  ('Food'),
                                  ('Beverages'),
                                  ('Household'),
                                  ('Personal Care'),
                                  ('Medicine'),
                                  ('Electronics'),
                                  ('Tools'),
                                  ('Other');

-- Insert default action types
INSERT INTO action_types (name) VALUES
                                    ('login'),
                                    ('logout'),
                                    ('create_item'),
                                    ('update_item'),
                                    ('delete_item'),
                                    ('create_storage'),
                                    ('update_storage'),
                                    ('delete_storage'),
                                    ('share_storage'),
                                    ('unshare_storage'),
                                    ('create_shopping_item'),
                                    ('update_shopping_item'),
                                    ('delete_shopping_item'),
                                    ('user_created'),
                                    ('user_updated'),
                                    ('user_deleted'),
                                    ('password_changed'),
                                    ('email_verified'),
                                    ('phone_verified'),
                                    ('two_factor_enabled'),
                                    ('two_factor_disabled');