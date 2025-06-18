BEGIN;

INSERT INTO categories (name) VALUES
  ('목걸이'),('반지'),('팔찌'),('귀걸이'),('금'),('은');
  
-- 2.2) 옵션 종류
INSERT INTO options (name) VALUES
  ('Size'),('Color'),('Purity');


-- 2.3) 옵션 값
INSERT INTO option_values (option_id, value)
SELECT o.id, v
FROM options o
JOIN (SELECT '15호' AS v UNION ALL SELECT '16호' UNION ALL SELECT '17호' UNION ALL SELECT '18호') AS vals ON o.name = 'Size'
UNION ALL
SELECT o.id, v
FROM options o
JOIN (SELECT '옐로우' AS v UNION ALL SELECT '실버' UNION ALL SELECT '로즈골드') AS vals ON o.name = 'Color'
UNION ALL
SELECT o.id, v
FROM options o
JOIN (SELECT '14k' AS v UNION ALL SELECT '18k' UNION ALL SELECT '24k') AS vals ON o.name = 'Purity';


INSERT INTO category_options (category_id, option_id)
SELECT c.id, o.id
FROM (
  SELECT '금','Purity' UNION ALL
  SELECT '금','Color' UNION ALL
  SELECT '은','Color' UNION ALL
  SELECT '목걸이','Size'   UNION ALL
  SELECT '반지','Size'     UNION ALL
  SELECT '팔찌','Size' 
) AS m(category_name, option_name)
JOIN categories c ON c.name = m.category_name
JOIN options    o ON o.name = m.option_name;



DELIMITER $$
CREATE PROCEDURE add_product_with_mappings(
  IN _name VARCHAR(100), 
  IN _price DECIMAL(10,2),
  IN _cats TEXT,      -- '금목걸이,목걸이' 식
  IN _opts TEXT       -- 'Size:16호,Size:17호,Size:18호,Color:옐로우,Color:그린,Color:로즈골드,Purity:18k,Purity:24k' 식
)
BEGIN
  INSERT INTO products (name, price) VALUES (_name, _price);
  SET @pid = LAST_INSERT_ID();

  -- 카테고리 매핑
  INSERT INTO product_categories (product_id, category_id)
  SELECT @pid, c.id
  FROM categories c
  WHERE FIND_IN_SET(c.name, _cats);

  -- 옵션 매핑
  INSERT INTO product_option_values (product_id, option_value_id)
  SELECT @pid, ov.id
  FROM option_values ov
  JOIN options o ON ov.option_id = o.id
  WHERE FIND_IN_SET(CONCAT(o.name, ':', ov.value), _opts);
END$$
DELIMITER ;



COMMIT;