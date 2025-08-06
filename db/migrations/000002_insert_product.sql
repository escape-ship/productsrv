BEGIN;


INSERT INTO products.products (id, name, description, base_price, image_url, category, options) VALUES
-- 금반지_모든옵션
('00000000-0000-0000-0000-000000000001', '금반지_모든옵션', '14k, 18k, 순금, 다양한 색상 가능', 300000,
'https://search.pstatic.net/common/?src=https%3A%2F%2Fshopping-phinf.pstatic.net%2Fmain_8881040%2F88810407822.2.jpg&type=f372_372',
'반지',
'{
  "metal": {
    "label": "금속 종류",
    "values": [
      {"name": "14k", "extra_price": 0},
      {"name": "18k", "extra_price": 30000},
      {"name": "순금", "extra_price": 100000}
    ]
  },
  "color": {
    "label": "색상",
    "dependency": "metal",
    "values": {
      "14k": ["yellow", "white"],
      "18k": ["yellow", "pink"],
      "순금": ["yellow"]
    }
  },
  "size": {
    "label": "사이즈",
    "values": ["10호", "11호", "12호", "13호", "14호", "15호"]
  }
}'::jsonb),

-- 금반지_14_18k만
('00000000-0000-0000-0000-000000000002', '금반지_14_18k만', '14k, 18k 전용 옵션', 280000,
'https://search.pstatic.net/common/?src=https%3A%2F%2Fsearchad-phinf.pstatic.net%2FMjAyNDA2MjRfMTkg%2FMDAxNzE5MjIwNzUwMDYx.7GUZmuFfKg32sDYNXf_Zr9b8rNcmwdHCSAy8Y9ZDkx8g.iqQXl_jU5mPNk46vKom82SzeBHYy6E8tNnhV-mFClC4g.PNG%2F1185960-1f77e095-962f-4e52-a060-52706dd96b65.png&type=f372_372',
'반지',
'{
  "metal": {
    "label": "금속 종류",
    "values": [
      {"name": "14k", "extra_price": 0},
      {"name": "18k", "extra_price": 30000}
    ]
  },
  "color": {
    "label": "색상",
    "dependency": "metal",
    "values": {
      "14k": ["yellow", "white"],
      "18k": ["yellow", "pink"]
    }
  },
  "size": {
    "label": "사이즈",
    "values": ["10호", "11호", "12호", "13호", "14호", "15호"]
  }
}'::jsonb),

-- 은반지_모든옵션
('00000000-0000-0000-0000-000000000003', '은반지_모든옵션', '은 전용, 색상 silver', 250000,
'https://search.pstatic.net/common/?src=https%3A%2F%2Fshopping-phinf.pstatic.net%2Fmain_8750081%2F87500811426.1.jpg&type=f372_372',

'반지',
'{
  "metal": {
    "label": "금속 종류",
    "values": [
      {"name": "silver", "extra_price": 0}
    ]
  },
  "color": {
    "label": "색상",
    "values": ["silver"]
  },
  "size": {
    "label": "사이즈",
    "values": ["10호", "11호", "12호", "13호", "14호"]
  }
}'::jsonb),

-- 금팔찌_모든옵션
('00000000-0000-0000-0000-000000000004', '금팔찌_모든옵션', '금 팔찌 다양한 옵션', 350000,
'https://search.pstatic.net/common/?src=https%3A%2F%2Fshopping-phinf.pstatic.net%2Fmain_8559585%2F85595859532.jpg&type=f372_372',
'팔찌',
'{
  "metal": {
    "label": "금속 종류",
    "values": [
      {"name": "14k", "extra_price": 0},
      {"name": "18k", "extra_price": 40000}
    ]
  },
  "color": {
    "label": "색상",
    "values": ["yellow", "pink"]
  },
  "length": {
    "label": "길이",
    "values": ["16cm", "18cm", "20cm"]
  }
}'::jsonb),

-- 금목걸이_모든옵션
('00000000-0000-0000-0000-000000000005', '금목걸이_모든옵션', '금 목걸이 다양한 옵션', 400000,
'https://search.pstatic.net/common/?src=http%3A%2F%2Fblogfiles.naver.net%2FMjAyMzA4MjVfMTQx%2FMDAxNjkyOTU4ODg3NDUx.cm1leBN3YaRek1TprZ4JrDJLJMyvi4rbOmnXLw7-9qkg.WFsJptkCxJ1PhFrcj6XYxqQvRa3xP1kOcTKf8kK3wpcg.JPEG.dterra%2F20230814-DSC03321-23-08-14.jpg&type=sc960_832',
'목걸이',
'{
  "metal": {
    "label": "금속 종류",
    "values": [
      {"name": "14k", "extra_price": 0},
      {"name": "18k", "extra_price": 50000}
    ]
  },
  "color": {
    "label": "색상",
    "values": ["yellow", "white"]
  },
  "length": {
    "label": "길이",
    "values": ["40cm", "45cm", "50cm"]
  }
}'::jsonb),

-- 은귀걸이_모든옵션
('00000000-0000-0000-0000-000000000006', '은귀걸이_모든옵션', '은 귀걸이', 270000,
'https://shop-phinf.pstatic.net/20240512_13/1715493682458iHtBd_JPEG/116629578194674771_1416493851.jpg?type=m510',
'귀걸이',
'{
  "metal": {
    "label": "금속 종류",
    "values": [
      {"name": "silver", "extra_price": 0}
    ]
  },
  "color": {
    "label": "색상",
    "values": ["silver"]
  }
}'::jsonb);

COMMIT;
