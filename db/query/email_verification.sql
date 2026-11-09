-- name: CreateEmailVerification :one
INSERT INTO email_verifications (email, code, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetEmailVerificationByEmailAndCode :one
SELECT * FROM email_verifications
WHERE email = $1
  AND code = $2
  AND used_at IS NULL
  AND expires_at > NOW()
ORDER BY created_at DESC
LIMIT 1;

-- name: GetVerifiedEmailVerification :one
SELECT * FROM email_verifications
WHERE email = $1
  AND verified = true
  AND used_at IS NULL
  AND expires_at > NOW()
ORDER BY created_at DESC
LIMIT 1;

-- name: MarkEmailVerificationVerified :exec
UPDATE email_verifications
SET verified = true
WHERE id = $1;

-- name: MarkEmailVerificationUsed :exec
UPDATE email_verifications
SET used_at = NOW()
WHERE id = $1;

-- name: DeleteExpiredEmailVerifications :exec
DELETE FROM email_verifications
WHERE expires_at < NOW();