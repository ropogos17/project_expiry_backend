-- ============================================================
-- 000003_users_uuid.up.sql
-- Convert users.id and every FK referencing users.id from INTEGER
-- to UUID. This includes auth_identities and refresh_tokens created
-- by migration 000002.
--
-- The original integer IDs are retained in users.legacy_integer_id
-- temporarily so the down migration can restore them if no new users
-- have been created since this migration.
-- ============================================================

-- ----------------------------------------
-- Preserve the original integer IDs and create UUIDs
-- ----------------------------------------

ALTER TABLE users
    ADD COLUMN legacy_integer_id INTEGER;

UPDATE users
SET legacy_integer_id = id;

ALTER TABLE users
    ALTER COLUMN legacy_integer_id DROP NOT NULL;

ALTER TABLE users
    ADD CONSTRAINT uq_users_legacy_integer_id
    UNIQUE (legacy_integer_id);

ALTER TABLE users
    ADD COLUMN id_uuid UUID;

UPDATE users
SET id_uuid = gen_random_uuid();

ALTER TABLE users
    ALTER COLUMN id_uuid SET NOT NULL;

ALTER TABLE users
    ALTER COLUMN id_uuid SET DEFAULT gen_random_uuid();

-- ----------------------------------------
-- Add and populate UUID shadow columns on every users FK
-- ----------------------------------------

ALTER TABLE auth_identities
    ADD COLUMN user_id_uuid UUID;

ALTER TABLE refresh_tokens
    ADD COLUMN user_id_uuid UUID;

ALTER TABLE storages
    ADD COLUMN owner_uuid UUID;

ALTER TABLE items
    ADD COLUMN last_modified_by_uuid UUID;

ALTER TABLE shares
    ADD COLUMN source_user_id_uuid UUID,
    ADD COLUMN target_user_id_uuid UUID;

ALTER TABLE announcements
    ADD COLUMN created_by_uuid UUID;

ALTER TABLE logs
    ADD COLUMN user_id_uuid UUID;

UPDATE auth_identities AS child
SET user_id_uuid = parent.id_uuid
FROM users AS parent
WHERE child.user_id = parent.legacy_integer_id;

UPDATE refresh_tokens AS child
SET user_id_uuid = parent.id_uuid
FROM users AS parent
WHERE child.user_id = parent.legacy_integer_id;

UPDATE storages AS child
SET owner_uuid = parent.id_uuid
FROM users AS parent
WHERE child.owner = parent.legacy_integer_id;

UPDATE items AS child
SET last_modified_by_uuid = parent.id_uuid
FROM users AS parent
WHERE child.last_modified_by = parent.legacy_integer_id;

UPDATE shares AS child
SET source_user_id_uuid = parent.id_uuid
FROM users AS parent
WHERE child.source_user_id = parent.legacy_integer_id;

UPDATE shares AS child
SET target_user_id_uuid = parent.id_uuid
FROM users AS parent
WHERE child.target_user_id = parent.legacy_integer_id;

UPDATE announcements AS child
SET created_by_uuid = parent.id_uuid
FROM users AS parent
WHERE child.created_by = parent.legacy_integer_id;

UPDATE logs AS child
SET user_id_uuid = parent.id_uuid
FROM users AS parent
WHERE child.user_id = parent.legacy_integer_id;

-- All original FK columns are NOT NULL. Enforce the same rule on the
-- populated UUID shadow columns before replacing the original columns.
ALTER TABLE auth_identities
    ALTER COLUMN user_id_uuid SET NOT NULL;

ALTER TABLE refresh_tokens
    ALTER COLUMN user_id_uuid SET NOT NULL;

ALTER TABLE storages
    ALTER COLUMN owner_uuid SET NOT NULL;

ALTER TABLE items
    ALTER COLUMN last_modified_by_uuid SET NOT NULL;

ALTER TABLE shares
    ALTER COLUMN source_user_id_uuid SET NOT NULL,
    ALTER COLUMN target_user_id_uuid SET NOT NULL;

ALTER TABLE announcements
    ALTER COLUMN created_by_uuid SET NOT NULL;

ALTER TABLE logs
    ALTER COLUMN user_id_uuid SET NOT NULL;

-- ----------------------------------------
-- Drop the old integer foreign keys
-- ----------------------------------------

ALTER TABLE auth_identities
    DROP CONSTRAINT auth_identities_user_id_fkey;

ALTER TABLE refresh_tokens
    DROP CONSTRAINT refresh_tokens_user_id_fkey;

ALTER TABLE storages
    DROP CONSTRAINT storages_owner_fkey;

ALTER TABLE items
    DROP CONSTRAINT items_last_modified_by_fkey;

ALTER TABLE shares
    DROP CONSTRAINT shares_source_user_id_fkey,
    DROP CONSTRAINT shares_target_user_id_fkey,
    DROP CONSTRAINT shares_check,
    DROP CONSTRAINT shares_storage_id_target_user_id_key;

ALTER TABLE announcements
    DROP CONSTRAINT announcements_created_by_fkey;

ALTER TABLE logs
    DROP CONSTRAINT logs_user_id_fkey;

-- ----------------------------------------
-- Replace the child integer FK columns with UUID columns
-- ----------------------------------------

ALTER TABLE auth_identities
    DROP COLUMN user_id;
ALTER TABLE auth_identities
    RENAME COLUMN user_id_uuid TO user_id;

ALTER TABLE refresh_tokens
    DROP COLUMN user_id;
ALTER TABLE refresh_tokens
    RENAME COLUMN user_id_uuid TO user_id;

ALTER TABLE storages
    DROP COLUMN owner;
ALTER TABLE storages
    RENAME COLUMN owner_uuid TO owner;

ALTER TABLE items
    DROP COLUMN last_modified_by;
ALTER TABLE items
    RENAME COLUMN last_modified_by_uuid TO last_modified_by;

ALTER TABLE shares
    DROP COLUMN source_user_id,
    DROP COLUMN target_user_id;
ALTER TABLE shares
    RENAME COLUMN source_user_id_uuid TO source_user_id;
ALTER TABLE shares
    RENAME COLUMN target_user_id_uuid TO target_user_id;

ALTER TABLE announcements
    DROP COLUMN created_by;
ALTER TABLE announcements
    RENAME COLUMN created_by_uuid TO created_by;

ALTER TABLE logs
    DROP COLUMN user_id;
ALTER TABLE logs
    RENAME COLUMN user_id_uuid TO user_id;

-- ----------------------------------------
-- Replace users.id INTEGER with users.id UUID
-- ----------------------------------------

ALTER TABLE users
    DROP CONSTRAINT users_pkey;

ALTER TABLE users
    ALTER COLUMN id DROP IDENTITY IF EXISTS;

ALTER TABLE users
    DROP COLUMN id;

ALTER TABLE users
    RENAME COLUMN id_uuid TO id;

ALTER TABLE users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

-- ----------------------------------------
-- Recreate foreign keys against the new UUID primary key
-- ----------------------------------------

ALTER TABLE auth_identities
    ADD CONSTRAINT auth_identities_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE refresh_tokens
    ADD CONSTRAINT refresh_tokens_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE storages
    ADD CONSTRAINT storages_owner_fkey
    FOREIGN KEY (owner) REFERENCES users (id);

ALTER TABLE items
    ADD CONSTRAINT items_last_modified_by_fkey
    FOREIGN KEY (last_modified_by) REFERENCES users (id);

ALTER TABLE shares
    ADD CONSTRAINT shares_source_user_id_fkey
        FOREIGN KEY (source_user_id) REFERENCES users (id),
    ADD CONSTRAINT shares_target_user_id_fkey
        FOREIGN KEY (target_user_id) REFERENCES users (id),
    ADD CONSTRAINT shares_storage_id_target_user_id_key
        UNIQUE (storage_id, target_user_id),
    ADD CONSTRAINT shares_check
        CHECK (source_user_id <> target_user_id);

ALTER TABLE announcements
    ADD CONSTRAINT announcements_created_by_fkey
    FOREIGN KEY (created_by) REFERENCES users (id);

ALTER TABLE logs
    ADD CONSTRAINT logs_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users (id);

-- ----------------------------------------
-- Recreate indexes whose columns were replaced
-- ----------------------------------------

CREATE INDEX idx_auth_identities_user_id
    ON auth_identities (user_id);

CREATE INDEX idx_refresh_tokens_user_id
    ON refresh_tokens (user_id);

CREATE INDEX idx_storages_owner
    ON storages (owner);

CREATE INDEX idx_shares_source_user_id
    ON shares (source_user_id);

CREATE INDEX idx_shares_target_user_id
    ON shares (target_user_id);

CREATE INDEX idx_announcements_created_by
    ON announcements (created_by);

CREATE INDEX idx_logs_user_id
    ON logs (user_id);

-- users.legacy_integer_id is intentionally retained for rollback.
-- New users created after this migration receive UUIDs and have
-- legacy_integer_id = NULL. That makes a later DOWN fail safely rather
-- than silently losing or changing user IDs.
