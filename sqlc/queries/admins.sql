-- name: GetAdmin :one
SELECT id, username FROM admins
WHERE id = $1;

-- name: ListAdmins :many
SELECT id, username FROM admins
ORDER BY id
LIMIT $1
OFFSET $2;
