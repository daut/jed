-- name: SaveToken :one
INSERT INTO tokens (hash, admin_id, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetToken :one
SELECT * FROM tokens
WHERE hash = $1;

-- name: DeleteTokens :exec
DELETE FROM tokens
WHERE admin_id = $1;
