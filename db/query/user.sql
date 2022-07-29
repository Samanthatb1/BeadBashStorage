-- name: CreateUser :one
INSERT INTO users (
  full_name,
  username,
  total_orders
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET total_orders = $2
WHERE id = $1
RETURNING *;;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username = $1;