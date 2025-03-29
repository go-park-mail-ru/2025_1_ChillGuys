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