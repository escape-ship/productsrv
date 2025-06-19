-- 카테고리: 소재(1:순금 2: 금, 3: 은), 종류(4: 목걸이, 5: 팔찌, 6: 반지, 7: 귀걸이)
INSERT INTO products.categories (id, name) VALUES
  (1, '순금'),
  (2, '금'),
  (3, '은'),
  (4, '목걸이'),
  (5, '팔찌'),
  (6, '반지'),
  (7, '귀걸이')
ON CONFLICT (id) DO NOTHING;