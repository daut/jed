-- name: GetAdmin :one
SELECT * FROM admins
WHERE username = $1;

-- name: GetAdminByID :one
SELECT * FROM admins
WHERE id = $1;

-- name: ListAdmins :many
SELECT id, username FROM admins
ORDER BY id
LIMIT $1
OFFSET $2;
