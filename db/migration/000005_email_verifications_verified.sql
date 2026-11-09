-- ============================================================
-- 000005_email_verifications_verified.sql
-- ============================================================

ALTER TABLE email_verifications
    ADD COLUMN verified BOOLEAN NOT NULL DEFAULT FALSE;
