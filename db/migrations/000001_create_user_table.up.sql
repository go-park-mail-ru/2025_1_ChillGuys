CREATE TYPE user_role AS ENUM ('seller', 'buyer', 'admin');


CREATE TABLE "user"
(
    user_id       UUID PRIMARY KEY,
    email         TEXT NOT NULL UNIQUE,
    phone_number  TEXT UNIQUE,
    password_hash TEXT NOT NULL,
    name          TEXT,
    surname       TEXT,
    image_url     TEXT,
    role          user_role,
    version       INTEGER
);

CREATE TABLE IF NOT EXISTS product (
    id                  UUID PRIMARY KEY,
    seller_id           UUID REFERENCES "user" (id) ON DELETE CASCADE,
    name                TEXT NOT NULL,
    preview_image_url   TEXT,
    description         TEXT,
    status              product_status NOT NULL,
    price               INT CHECK (price >= 0) NOT NULL,
    quantity            INT CHECK (quantity >= 0) NOT NULL,
    updated_at          TIMESTAMPTZ DEFAULT now(),
    rating              INT CHECK (rating BETWEEN 1 AND 5)
);