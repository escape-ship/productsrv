-- name: GetProducts :many
SELECT p.id, p.name, p.price, p.image_url, p.created_at, p.updated_at,
       COALESCE(json_agg(DISTINCT c.name) FILTER (WHERE c.id IS NOT NULL), '[]') AS categories,
       COALESCE(SUM(i.stock_quantity), 0) AS inventory
FROM products.product p
LEFT JOIN products.products_categories_relations pcr ON p.id = pcr.product_id
LEFT JOIN products.categories c ON pcr.category_id = c.id
LEFT JOIN products.inventories i ON p.id = i.product_id
GROUP BY p.id;

-- name: GetProductByName :one
SELECT p.id, p.name, p.price, p.image_url, p.created_at, p.updated_at,
       COALESCE(json_agg(DISTINCT c.name) FILTER (WHERE c.id IS NOT NULL), '[]') AS categories,
       COALESCE(SUM(i.stock_quantity), 0) AS inventory
FROM products.product p
LEFT JOIN products.products_categories_relations pcr ON p.id = pcr.product_id
LEFT JOIN products.categories c ON pcr.category_id = c.id
LEFT JOIN products.inventories i ON p.id = i.product_id
WHERE p.name = $1
GROUP BY p.id;

-- name: PostProducts :exec
INSERT INTO products.product (id, name, price, image_url)
VALUES ($1, $2, $3, $4);

-- name: GetInventoriesByProductID :many
SELECT id, product_id, product_option_id, stock_quantity
FROM products.inventories
WHERE product_id = $1;

