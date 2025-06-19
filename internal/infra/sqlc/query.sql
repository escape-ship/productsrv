-- name: GetProducts :many
SELECT id, name, price, image_url, created_at, updated_at
FROM products.product;

-- name: GetProductByName :one
SELECT id, name, price, image_url, created_at, updated_at
FROM products.product
WHERE name = $1;

-- name: GetProductByID :one
SELECT id, name, price, image_url, created_at, updated_at
FROM products.product
WHERE id = $1;

-- name: PostProducts :exec
INSERT INTO products.product (id, name, price, image_url)
VALUES ($1, $2, $3, $4);

-- name: GetCategories :many
SELECT id, name
FROM products.categories;

-- name: GetCategoryByID :one
SELECT id, name
FROM products.categories
WHERE id = $1;

-- name: GetCategoriesByProductID :many
SELECT c.id, c.name
FROM products.categories c
JOIN products.product_categories pc ON c.id = pc.category_id
WHERE pc.product_id = $1;

-- name: HasProductCustomOptionValues :one
SELECT EXISTS (
  SELECT 1
  FROM products.product_option_values
  WHERE product_id = $1
) AS has_custom;

-- name: GetProductCustomOptionValues :many
SELECT
  o.id AS option_id,
  o.name AS option_name,
  ov.id AS value_id,
  ov.value AS value
FROM products.product_option_values pov
JOIN products.option_values ov ON pov.option_value_id = ov.id
JOIN products.options o ON pov.option_id = o.id
WHERE pov.product_id = $1
ORDER BY o.id, ov.id;

-- name: GetProductDefaultOptionValues :many
SELECT
  o.id AS option_id,
  o.name AS option_name,
  ov.id AS value_id,
  ov.value AS value
FROM products.product_categories pc
JOIN products.category_options co ON pc.category_id = co.category_id
JOIN products.options o ON co.option_id = o.id
JOIN products.option_values ov ON ov.option_id = o.id
WHERE pc.product_id = $1
ORDER BY o.id, ov.id;
