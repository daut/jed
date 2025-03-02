-- name: CreateOrder :one
INSERT INTO orders (first_name, last_name, email, phone, address, city)
VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: CreateOrderProduct :many
INSERT INTO order_product (order_id, product_id, quantity)
VALUES (
  $1, $2, $3
)
RETURNING *;
