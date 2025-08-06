-- name: GetProducts :many
SELECT id, name, base_price AS price, image_url, category, options, created_at, updated_at
FROM products.products;

-- name: GetProductByName :one
SELECT id, name, base_price AS price, image_url, category, options, created_at, updated_at
FROM products.products
WHERE name = $1;

-- name: GetProductByID :one
SELECT id, name, description, base_price AS price, image_url, category, options, created_at, updated_at
FROM products.products
WHERE id = $1;

-- name: PostProducts :exec
INSERT INTO products.products (name, description, base_price, image_url, category, options)
VALUES ($1, $2, $3, $4, $5, $6);
