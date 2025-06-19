-- 카테고리-옵션 매핑
-- 금(1) → 함량(1), 색상(3)
INSERT INTO products.category_options (category_id, option_id) VALUES
  (2, 1), (2, 3);
-- 반지(5) → 사이즈(2)
INSERT INTO products.category_options (category_id, option_id) VALUES
  (5, 2);
