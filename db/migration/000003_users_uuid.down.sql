-- ============================================================
-- 000003_users_uuid.down.sql
-- Restore users.id and all its referencing FKs to INTEGER.
--
-- This rollback is guarded: it is only possible while every user still
-- has an original integer ID in users.legacy_integer_id. If users were
-- created after 000003, the migration stops before making any changes.
-- ============================================================

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM users
        WHERE legacy_integer_id IS NULL
    ) THEN
        RAISE EXCEPTION
            'Cannot roll back 000003: one or more users were created after the UUID migration and have no legacy integer ID.';
    END IF;
END;
$$;

-- ----------------------------------------
-- Add and populate INTEGER shadow columns on every users FK
-- ----------------------------------------

ALTER TABLE auth_identities
    ADD COLUMN user_id_integer INTEGER;

ALTER TABLE refresh_tokens
    ADD COLUMN user_id_integer INTEGER;

ALTER TABLE storages
    ADD COLUMN owner_integer INTEGER;

ALTER TABLE items
    ADD COLUMN last_modified_by_integer INTEGER;

ALTER TABLE shares
    ADD COLUMN source_user_id_integer INTEGER,
    ADD COLUMN target_user_id_integer INTEGER;

ALTER TABLE announcements
    ADD COLUMN created_by_integer INTEGER;

ALTER TABLE logs
    ADD COLUMN user_id_integer INTEGER;

UPDATE auth_identities AS child
SET user_id_integer = parent.legacy_integer_id
FROM users AS parent
WHERE child.user_id = parent.id;

UPDATE refresh_tokens AS child
SET user_id_integer = parent.legacy_integer_id
FROM users AS parent
WHERE child.user_id = parent.id;

UPDATE storages AS child
SET owner_integer = parent.legacy_integer_id
FROM users AS parent
WHERE child.owner = parent.id;

UPDATE items AS child
SET last_modified_by_integer = parent.legacy_integer_id
FROM users AS parent
WHERE child.last_modified_by = parent.id;

UPDATE shares AS child
SET source_user_id_integer = parent.legacy_integer_id
FROM users AS parent
WHERE child.source_user_id = parent.id;

UPDATE shares AS child
SET target_user_id_integer = parent.legacy_integer_id
FROM users AS parent
WHERE child.target_user_id = parent.id;

UPDATE announcements AS child
SET created_by_integer = parent.legacy_integer_id
FROM users AS parent
WHERE child.created_by = parent.id;

UPDATE logs AS child
SET user_id_integer = parent.legacy_integer_id
FROM users AS parent
WHERE child.user_id = parent.id;

ALTER TABLE auth_identities
    ALTER COLUMN user_id_integer SET NOT NULL;

ALTER TABLE refresh_tokens
    ALTER COLUMN user_id_integer SET NOT NULL;

ALTER TABLE storages
    ALTER COLUMN owner_integer SET NOT NULL;

ALTER TABLE items
    ALTER COLUMN last_modified_by_integer SET NOT NULL;

ALTER TABLE shares
    ALTER COLUMN source_user_id_integer SET NOT NULL,
    ALTER COLUMN target_user_id_integer SET NOT NULL;

ALTER TABLE announcements
    ALTER COLUMN created_by_integer SET NOT NULL;

ALTER TABLE logs
    ALTER COLUMN user_id_integer SET NOT NULL;

-- ----------------------------------------
-- Drop the UUID foreign keys
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
-- Replace the child UUID FK columns with INTEGER columns
-- ----------------------------------------

ALTER TABLE auth_identities
    DROP COLUMN user_id;
ALTER TABLE auth_identities
    RENAME COLUMN user_id_integer TO user_id;

ALTER TABLE refresh_tokens
    DROP COLUMN user_id;
ALTER TABLE refresh_tokens
    RENAME COLUMN user_id_integer TO user_id;

ALTER TABLE storages
    DROP COLUMN owner;
ALTER TABLE storages
    RENAME COLUMN owner_integer TO owner;

ALTER TABLE items
    DROP COLUMN last_modified_by;
ALTER TABLE items
    RENAME COLUMN last_modified_by_integer TO last_modified_by;

ALTER TABLE shares
    DROP COLUMN source_user_id,
    DROP COLUMN target_user_id;
ALTER TABLE shares
    RENAME COLUMN source_user_id_integer TO source_user_id;
ALTER TABLE shares
    RENAME COLUMN target_user_id_integer TO target_user_id;

ALTER TABLE announcements
    DROP COLUMN created_by;
ALTER TABLE announcements
    RENAME COLUMN created_by_integer TO created_by;

ALTER TABLE logs
    DROP COLUMN user_id;
ALTER TABLE logs
    RENAME COLUMN user_id_integer TO user_id;

-- ----------------------------------------
-- Restore users.id INTEGER using the preserved original IDs
-- ----------------------------------------

ALTER TABLE users
    DROP CONSTRAINT users_pkey;

ALTER TABLE users
    ADD COLUMN id_integer INTEGER;

UPDATE users
SET id_integer = legacy_integer_id;

ALTER TABLE users
    ALTER COLUMN id_integer SET NOT NULL;

ALTER TABLE users
    DROP COLUMN id;

ALTER TABLE users
    RENAME COLUMN id_integer TO id;

ALTER TABLE users
    ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY;

ALTER TABLE users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

-- Set the newly created identity sequence above the restored maximum ID.
SELECT setval(
    pg_get_serial_sequence('users', 'id'),
    COALESCE((SELECT MAX(id) FROM users), 1),
    EXISTS (SELECT 1 FROM users)
);

-- ----------------------------------------
-- Recreate foreign keys against the restored INTEGER primary key
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

-- Remove the temporary rollback mapping after all IDs are restored.
ALTER TABLE users
    DROP CONSTRAINT uq_users_legacy_integer_id,
    DROP COLUMN legacy_integer_id;
