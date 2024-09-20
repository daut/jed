-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1;

-- name: GetProducts :many
SELECT * FROM products
ORDER BY id
LIMIT $1
OFFSET $2;
