-- ============================================================
-- 000002_auth_upgrade.down.sql
-- Roll back the external-identity / hashed-refresh-token changes.
--
-- NOTE: This rollback is intentionally guarded. It cannot restore the
-- old NOT NULL email/password rules if passwordless users now exist.
-- ============================================================

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM users
        WHERE email IS NULL OR password_hash IS NULL
    ) THEN
        RAISE EXCEPTION
            'Cannot roll back 000002: users with NULL email or password_hash exist. '
            'Resolve or migrate those accounts before running the down migration.';
    END IF;
END;
$$;

-- External identity links are removed by this rollback.
DROP TABLE auth_identities;

ALTER TABLE users
    DROP CONSTRAINT chk_users_password_requires_email;

-- Refresh sessions created with token_hash are invalidated on rollback.
DELETE FROM refresh_tokens;

ALTER TABLE refresh_tokens
    DROP CONSTRAINT uq_refresh_tokens_token_hash,
    DROP CONSTRAINT chk_refresh_tokens_token_hash_format,
    DROP COLUMN token_hash;

ALTER TABLE users
    RENAME COLUMN password_hash TO password;

ALTER TABLE users
    ALTER COLUMN email SET NOT NULL,
    ALTER COLUMN password SET NOT NULL;
