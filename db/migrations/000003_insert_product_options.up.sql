-- 옵션: 함량(1), 사이즈(2), 색상(3)
INSERT INTO products.options (id, name) VALUES
  (1, '함량'),
  (2, '사이즈'),
  (3, '색상');

-- 옵션값: 함량
INSERT INTO products.option_values (id, option_id, value) VALUES
  (1, 1, '14k'),
  (2, 1, '18k');

-- 옵션값: 사이즈 (10호 ~ 20호)
INSERT INTO products.option_values (id, option_id, value) VALUES
  (4, 2, '10호'),
  (5, 2, '11호'),
  (6, 2, '12호'),
  (7, 2, '13호'),
  (8, 2, '14호'),
  (9, 2, '15호'),
  (10, 2, '16호'),
  (11, 2, '17호'),
  (12, 2, '18호'),
  (13, 2, '19호'),
  (14, 2, '20호');

-- 옵션값: 색상
INSERT INTO products.option_values (id, option_id, value) VALUES
  (15, 3, 'yellow'),
  (16, 3, 'pink'),
  (17, 3, 'silver');