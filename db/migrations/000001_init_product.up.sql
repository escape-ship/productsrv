BEGIN;

CREATE SCHEMA IF NOT EXISTS products;

CREATE TABLE products.product (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    price BIGINT NOT NULL,
    image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products.product_options (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products.product(id) ON DELETE CASCADE,
    option TEXT NOT NULL,
    value TEXT NOT NULL
);

CREATE TABLE products.inventories (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products.product(id) ON DELETE CASCADE,
    product_option_id UUID NOT NULL REFERENCES products.product_options(id) ON DELETE CASCADE,
    stock_quantity INT NOT NULL
);

CREATE TABLE products.categories (
    id UUID PRIMARY KEY,
    parent_id UUID REFERENCES products.categories(id) ON DELETE SET NULL,
    name TEXT NOT NULL
);

CREATE TABLE products.products_categories_relations (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products.product(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES products.categories(id) ON DELETE CASCADE
);


COMMIT;