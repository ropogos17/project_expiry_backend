-- ============================================================
-- 000005_email_verifications_verified.down.sql
-- ============================================================

ALTER TABLE email_verifications
    DROP COLUMN verified;
