-- name: CreateOrder :one
INSERT INTO orders (
  account_id,
  username,
  full_name,
  purchase_amount,
  purchased_item,
  shipping_location,
  currency,
  date_ordered
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetOrderById :one
SELECT * FROM orders
WHERE order_id = $1 LIMIT 1;

-- name: ListOrdersByUsername :many
SELECT * FROM orders
WHERE username = $1;

-- name: ListAllOrders :many
SELECT * FROM orders
ORDER BY order_id
LIMIT $1
OFFSET $2;

-- name: UpdateOrder :one
UPDATE orders
SET purchase_amount = $2,
purchased_item = $3,
shipping_location = $4
WHERE order_id = $1
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE order_id = $1;

-- name: DeleteAllOrderFromUser :exec
DELETE FROM orders
WHERE username = $1;