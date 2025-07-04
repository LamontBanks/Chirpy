-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token, user_id, created_at, updated_at, expires_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
);


-- name: GetRefreshTokenInfo :one
SELECT token, user_id, created_at, updated_at, expires_at, revoked_at FROM refresh_tokens
WHERE token = $1;

-- name: SetRefreshTokenRevokeAtTime :exec
UPDATE refresh_tokens
SET revoked_at = $2, updated_at = $3
WHERE token = $1;
