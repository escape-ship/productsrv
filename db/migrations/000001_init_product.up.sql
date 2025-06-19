BEGIN;

CREATE SCHEMA IF NOT EXISTS products;

CREATE TABLE IF NOT EXISTS products.product (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    price BIGINT NOT NULL,
    image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products.categories (
    id    INT PRIMARY KEY,
    name  VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products.product_categories (
  product_id UUID,
  category_id INT,
  PRIMARY KEY (product_id, category_id),
  FOREIGN KEY (product_id)  REFERENCES products.product(id),
  FOREIGN KEY (category_id) REFERENCES products.categories(id)
);

CREATE TABLE IF NOT EXISTS products.options(
    id    INT PRIMARY KEY,
    name  VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS products.option_values (
  id        INT PRIMARY KEY,
  option_id INT NOT NULL,
  value     VARCHAR(50) NOT NULL,
  UNIQUE (option_id, value),
  FOREIGN KEY (option_id) REFERENCES products.options(id)
);

-- 카테고리별 허용 옵션 매핑
CREATE TABLE IF NOT EXISTS products.category_options (
  category_id INT NOT NULL,
  option_id   INT NOT NULL,
  PRIMARY KEY (category_id, option_id),
  FOREIGN KEY (category_id) REFERENCES products.categories(id),
  FOREIGN KEY (option_id)   REFERENCES products.options(id)
);


-- M:N 매핑: 상품 ↔ 옵션 값
CREATE TABLE IF NOT EXISTS products.product_option_values (
  product_id      UUID  NOT NULL,
  option_value_id INT NOT NULL,
  option_id INT NOT NULL,
  PRIMARY KEY (product_id, option_value_id),
  FOREIGN KEY (product_id)      REFERENCES products.product(id),
  FOREIGN KEY (option_value_id) REFERENCES products.option_values(id),
  FOREIGN KEY (option_id) REFERENCES products.options(id)
);

COMMIT;