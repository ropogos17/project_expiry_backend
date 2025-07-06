-- Create roles table
CREATE TABLE roles (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) UNIQUE NOT NULL
);

-- Create status table
CREATE TABLE status (
                        id SERIAL PRIMARY KEY,
                        name VARCHAR(255) UNIQUE NOT NULL
);

-- Create users table
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       username VARCHAR(255) UNIQUE NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       phone VARCHAR(255) DEFAULT '',
                       has_avatar BOOLEAN DEFAULT FALSE,
                       password VARCHAR(255) NOT NULL,
                       salt VARCHAR(255) NOT NULL,
                       email_verified BOOLEAN DEFAULT FALSE,
                       phone_verified BOOLEAN DEFAULT FALSE,
                       two_factor_enabled BOOLEAN DEFAULT FALSE,
                       two_factor_secret VARCHAR(255),
                       last_login TIMESTAMP,
                       login_attempts INTEGER DEFAULT 0,
                       locked_until TIMESTAMP,
                       password_changed_at TIMESTAMP NOT NULL,
                       role_id INTEGER NOT NULL REFERENCES roles(id),
                       is_superuser BOOLEAN DEFAULT FALSE,
                       status_id INTEGER NOT NULL REFERENCES status(id),
                       deleted_at TIMESTAMP,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP
);

-- Create categories table
CREATE TABLE categories (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(255) NOT NULL
);

-- Create storage table
CREATE TABLE storage (
                         id SERIAL PRIMARY KEY,
                         name VARCHAR(255) NOT NULL,
                         description TEXT,
                         location VARCHAR(255),
                         shelves INTEGER,
                         has_image BOOLEAN DEFAULT FALSE,
                         owner INTEGER NOT NULL REFERENCES users(id)
);

-- Create items table
CREATE TABLE items (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       category_id INTEGER NOT NULL REFERENCES categories(id),
                       description TEXT,
                       quantity TEXT,
                       has_image BOOLEAN DEFAULT FALSE,
                       is_deleted BOOLEAN DEFAULT FALSE,
                       stored_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP,
                       expires_at TIMESTAMP NOT NULL,
                       storage_id INTEGER NOT NULL REFERENCES storage(id),
                       shelf INTEGER,
                       last_modified INTEGER NOT NULL REFERENCES users(id)
);

-- Create announcements table
CREATE TABLE announcements (
                               id SERIAL PRIMARY KEY,
                               name VARCHAR(255) NOT NULL,
                               description TEXT,
                               external_link TEXT,
                               has_read BOOLEAN DEFAULT FALSE
);

-- Create shopping_list_items table
CREATE TABLE shopping_list_items (
                                     id SERIAL PRIMARY KEY,
                                     name VARCHAR(255) NOT NULL,
                                     description TEXT,
                                     quantity TEXT,
                                     is_checked BOOLEAN DEFAULT FALSE,
                                     is_deleted BOOLEAN DEFAULT FALSE,
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMP,
                                     last_modified INTEGER NOT NULL REFERENCES users(id)
);

-- Create shares table
CREATE TABLE shares (
                        id SERIAL PRIMARY KEY,
                        source_user_id INTEGER NOT NULL REFERENCES users(id),
                        storage_id INTEGER NOT NULL REFERENCES storage(id),
                        target_user_id INTEGER NOT NULL REFERENCES users(id),
                        can_share BOOLEAN DEFAULT FALSE,
                        shared_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create languages table
CREATE TABLE languages (
                           id SERIAL PRIMARY KEY,
                           en TEXT NOT NULL,
                           hu TEXT NOT NULL,
                           ro TEXT NOT NULL,
                           es TEXT NOT NULL,
                           zh TEXT NOT NULL,
                           fr TEXT NOT NULL,
                           de TEXT NOT NULL,
                           pt TEXT NOT NULL,
                           ru TEXT NOT NULL,
                           ua TEXT NOT NULL,
                           ar TEXT NOT NULL
);

-- Create action_types table
CREATE TABLE action_types (
                              id SERIAL PRIMARY KEY,
                              name VARCHAR(255) UNIQUE NOT NULL
);

-- Create logs table
CREATE TABLE logs (
                      id SERIAL PRIMARY KEY,
                      created_at TIMESTAMP DEFAULT NOW(),
                      action_type INTEGER REFERENCES action_types(id),
                      message TEXT,
                      data TEXT,
                      user_id INTEGER NOT NULL REFERENCES users(id)
);
