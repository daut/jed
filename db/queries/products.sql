-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1;

-- name: GetProducts :many
SELECT * FROM products
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CreateProduct :one
INSERT INTO products (name, description, price, inventory_count)
VALUES (
  $1, $2, $3,
  coalesce(sqlc.narg(inventory_count)::integer, 1)
)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET name = $1, description = $2, price = $3, inventory_count = $4
WHERE id = $5
RETURNING *;

-- name: DeleteProduct :one
DELETE FROM products
WHERE id = $1
RETURNING *;
