-- name: GetProduct :many
SELECT id, name, categories, price, inventory, imageUrl, created_at, updated_at
FROM products.product
WHERE $1 = 0 OR id = $1;

-- name: GetProductByName :one
SELECT id, name, categories, price, inventory, imageUrl, created_at, updated_at
FROM products.product
WHERE name = $1;

-- name: PostProducts :exec
INSERT INTO products.product (name, categories, price, inventory, imageUrl)
VALUES ($1, $2, $3, $4, $5);

