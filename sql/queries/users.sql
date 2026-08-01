-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (
    gen_random_uuid()
    ,NOW()
    ,NOW()
    ,$1
    ,$2
    ,false
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdatePassword :exec
UPDATE users
set hashed_password = $1
    , updated_at = NOW()
WHERE id = $2;

-- name: UpdatePasswordAndEmail :one
UPDATE users
SET hashed_password = $1
    , email = $2
    , updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: UpgradeToChirpyRed :exec
UPDATE users
SET is_chirpy_red = true
WHERE id = $1;