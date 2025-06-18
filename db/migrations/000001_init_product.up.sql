BEGIN;

CREATE SCHEMA products;

CREATE TABLE products.product (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    price BIGINT NOT NULL,
    image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products.categories (
    id    INT,
    name  VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE products.product_categories(
    product_id UUID,
    category_id INT,
    PRIMARY KEY (product_id, category_id),
    FOREIGN KEY (product_id)  REFERENCES products(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE products.options(
    id    INT PRIMARY KEY,
    name  VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE products.option_values (
  id        INT PRIMARY KEY,
  option_id INT NOT NULL,
  value     VARCHAR(50) NOT NULL,
  UNIQUE (option_id, value),
  FOREIGN KEY (option_id) REFERENCES options(id)
);

-- 카테고리별 허용 옵션 매핑
CREATE TABLE products.category_options (
  category_id INT NOT NULL,
  option_id   INT NOT NULL,
  PRIMARY KEY (category_id, option_id),
  FOREIGN KEY (category_id) REFERENCES categories(id),
  FOREIGN KEY (option_id)   REFERENCES options(id)
);


-- M:N 매핑: 상품 ↔ 옵션 값
CREATE TABLE products.product_option_values (
  product_id      INT NOT NULL,
  option_value_id INT NOT NULL,
  PRIMARY KEY (product_id, option_value_id),
  FOREIGN KEY (product_id)      REFERENCES products(id),
  FOREIGN KEY (option_value_id) REFERENCES option_values(id)
);

COMMIT;