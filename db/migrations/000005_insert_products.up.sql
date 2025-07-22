-- 상품 목록 추가
INSERT INTO products.product (id, name, price, image_url) VALUES
  ('00000000-0000-0000-0000-000000000001', '금반지_모든옵션', 300000, 'https://search.pstatic.net/common/?src=https%3A%2F%2Fshopping-phinf.pstatic.net%2Fmain_8881040%2F88810407822.2.jpg&type=f372_372'),
  ('00000000-0000-0000-0000-000000000002', '금반지_14_18k만', 280000, 'https://search.pstatic.net/common/?src=https%3A%2F%2Fsearchad-phinf.pstatic.net%2FMjAyNDA2MjRfMTkg%2FMDAxNzE5MjIwNzUwMDYx.7GUZmuFfKg32sDYNXf_Zr9b8rNcmwdHCSAy8Y9ZDkx8g.iqQXl_jU5mPNk46vKom82SzeBHYy6E8tNnhV-mFClC4g.PNG%2F1185960-1f77e095-962f-4e52-a060-52706dd96b65.png&type=f372_372'),
  ('00000000-0000-0000-0000-000000000003', '은반지_모든옵션', 250000, 'https://search.pstatic.net/common/?src=https%3A%2F%2Fshopping-phinf.pstatic.net%2Fmain_8750081%2F87500811426.1.jpg&type=f372_372'),
  ('00000000-0000-0000-0000-000000000004', '금팔찌_모든옵션', 350000, 'https://search.pstatic.net/common/?src=https%3A%2F%2Fshopping-phinf.pstatic.net%2Fmain_8559585%2F85595859532.jpg&type=f372_372'),
  ('00000000-0000-0000-0000-000000000005', '금목걸이_모든옵션', 400000, 'https://search.pstatic.net/common/?src=http%3A%2F%2Fblogfiles.naver.net%2FMjAyMzA4MjVfMTQx%2FMDAxNjkyOTU4ODg3NDUx.cm1leBN3YaRek1TprZ4JrDJLJMyvi4rbOmnXLw7-9qkg.WFsJptkCxJ1PhFrcj6XYxqQvRa3xP1kOcTKf8kK3wpcg.JPEG.dterra%2F20230814-DSC03321-23-08-14.jpg&type=sc960_832'),
  ('00000000-0000-0000-0000-000000000006', '은귀걸이_모든옵션', 270000, 'https://shop-phinf.pstatic.net/20240512_13/1715493682458iHtBd_JPEG/116629578194674771_1416493851.jpg?type=m510');

-- 상품별 카테고리 매핑
INSERT INTO products.product_categories (product_id, category_id) VALUES
  -- 1번 금반지 (금=2, 반지=6)
  ('00000000-0000-0000-0000-000000000001', 2),
  ('00000000-0000-0000-0000-000000000001', 6),

  -- 2번 금반지 14,18k만 (금=2, 반지=6)
  ('00000000-0000-0000-0000-000000000002', 2),
  ('00000000-0000-0000-0000-000000000002', 6),

  -- 3번 은반지 (은=3, 반지=6)
  ('00000000-0000-0000-0000-000000000003', 3),
  ('00000000-0000-0000-0000-000000000003', 6),

  -- 4번 금팔찌 (금=2, 팔찌=5)
  ('00000000-0000-0000-0000-000000000004', 2),
  ('00000000-0000-0000-0000-000000000004', 5),

  -- 5번 금목걸이 (금=2, 목걸이=4)
  ('00000000-0000-0000-0000-000000000005', 2),
  ('00000000-0000-0000-0000-000000000005', 4),

  -- 6번 은귀걸이 (은=3, 귀걸이=7)
  ('00000000-0000-0000-0000-000000000006', 3),
  ('00000000-0000-0000-0000-000000000006', 7);

-- 허용되는 옵션이 다른 옵션 추가
INSERT INTO products.product_option_values (product_id, option_value_id, option_id) VALUES
  ('00000000-0000-0000-0000-000000000002', 1, 1),  -- 14k
  ('00000000-0000-0000-0000-000000000002', 2, 1);  -- 18k