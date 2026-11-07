-- ============================================================
-- 000004_email_verifications.down.sql
-- ============================================================

DROP INDEX IF EXISTS idx_email_verifications_code;
DROP INDEX IF EXISTS idx_email_verifications_email;
DROP TABLE IF EXISTS email_verifications;