BEGIN;

-- 확장 기능: UUID 생성 함수
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 스키마 생성
CREATE SCHEMA IF NOT EXISTS products;

-- 상품 테이블
CREATE TABLE IF NOT EXISTS products.products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,                            -- 상품명
    description TEXT,                              -- 상세 설명
    base_price BIGINT NOT NULL,                   -- 기본 가격
    image_url TEXT,
    category TEXT,
    options JSONB,                                 -- 옵션 구조(JSON)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMIT;
