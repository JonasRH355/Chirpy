-- name: NewToken :exec
INSERT INTO refresh_tokens (
    token,
    created_at,
    updated_at
    ,user_id
    ,expires_at
) VALUES (
    $1
    ,NOW()
    ,NOW()
    ,$2
    ,NOW() + INTERVAL '60 days'
);

-- name: GetUserFromRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1
    AND NOW() < expires_at;

-- name: UpdateRefreshToken :exec
UPDATE refresh_tokens
SET 
    revoked_at = NOW(),
    updated_at = NOW() + INTERVAL '60 days'
WHERE token = $1;