-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1;

-- name: GetProducts :many
SELECT * FROM products
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreateProduct :one
INSERT INTO products (name, description, price)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products
WHERE id = $1
RETURNING *;
