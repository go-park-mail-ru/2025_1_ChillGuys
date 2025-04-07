-- Создание алиаса для схемы bazaar
CREATE SCHEMA IF NOT EXISTS bazaar;

-- Создание ENUM для статуса товара
CREATE TYPE bazaar.product_status AS ENUM (
    'pending', -- Ожидает
    'rejected', -- Отказано
    'approved' -- Одобрено
    );

-- Создание ENUM для статуса заказа
CREATE TYPE bazaar.order_status AS ENUM (
    'pending', -- Ожидает
    'placed', -- Оформлен
    'awaiting_confirmation', -- Ожидает подтверждения
    'being_prepared', -- Готовится
    'shipped', -- Отправлен
    'in_transit', -- В пути
    'delivered_to_pickup_point', -- Доставлен в пункт самовывоза
    'delivered', -- Доставлен
    'canceled', -- Отменен
    'awaiting_payment', -- Ожидает оплаты
    'paid', -- Оплачено
    'payment_failed', -- Платеж не удался
    'return_requested', -- Возврат запрашивается
    'return_processed', -- Возврат обработан
    'return_initiated', -- Возврат инициирован
    'return_completed', -- Возврат завершен
    'canceled_by_user', -- Отменен пользователем
    'canceled_by_seller', -- Отменен продавцом
    'canceled_due_to_payment_error' -- Отменен из-за ошибки платежа
    );

-- Адреса
CREATE TABLE IF NOT EXISTS bazaar.address
(
    id         UUID PRIMARY KEY,
    city       TEXT NOT NULL,
    street     TEXT NOT NULL,
    house      TEXT NOT NULL,
    apartment  TEXT,
    zip_code   TEXT NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS bazaar."user"
(
    id            UUID PRIMARY KEY,
    email         TEXT UNIQUE NOT NULL,
    phone_number  TEXT,
    password_hash TEXT        NOT NULL,
    name          TEXT        NOT NULL,
    surname       TEXT,
    image_url     TEXT,
    address_id    UUID        REFERENCES bazaar.address (id) ON DELETE SET NULL
);

-- Таблица для ПВЗ
CREATE TABLE IF NOT EXISTS bazaar.pickup_point
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    address_id UUID NOT NULL REFERENCES bazaar.address (id) ON DELETE SET NULL
);

-- Связующая таблица для ПВЗ и пользователей
CREATE TABLE IF NOT EXISTS bazaar.user_pickup_point
(
    id              UUID PRIMARY KEY,
    user_id         UUID REFERENCES bazaar."user" (id) ON DELETE CASCADE,
    pickup_point_id UUID REFERENCES bazaar.pickup_point (id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE bazaar.user_balance
(
    id         UUID PRIMARY KEY,
    user_id    UUID           NOT NULL REFERENCES bazaar."user" (id) ON DELETE CASCADE,
    balance    NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    UNIQUE (user_id)
);

-- Создание таблицы ролей
CREATE TABLE IF NOT EXISTS bazaar.role
(
    id   UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Создание таблицы user_role
CREATE TABLE IF NOT EXISTS bazaar.user_role
(
    id      UUID PRIMARY KEY,
    user_id UUID REFERENCES bazaar."user" (id) ON DELETE CASCADE,
    role_id UUID REFERENCES bazaar.role (id) ON DELETE CASCADE,
    UNIQUE (user_id, role_id)
);

-- Версии пользователя
CREATE TABLE IF NOT EXISTS bazaar.user_version
(
    id         UUID PRIMARY KEY,
    user_id    UUID REFERENCES bazaar."user" (id) ON DELETE CASCADE,
    version    INT         DEFAULT 1 CHECK (version > 0),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Товары
CREATE TABLE IF NOT EXISTS bazaar.product
(
    id                UUID PRIMARY KEY,
    seller_id         UUID REFERENCES bazaar."user" (id) ON DELETE CASCADE,
    name              TEXT                              NOT NULL,
    preview_image_url TEXT                               DEFAULT 'media/product-default',
    description       TEXT,
    status            bazaar.product_status             NOT NULL,
    price             NUMERIC(12, 2) CHECK (price >= 0) NOT NULL,
    quantity          INT CHECK (quantity >= 0)         NOT NULL,
    updated_at        TIMESTAMPTZ                        DEFAULT now(),
    rating            INT CHECK (rating BETWEEN 0 AND 5) DEFAULT 0,
    reviews_count     INT CHECK (reviews_count >= 0)     DEFAULT 0 --trigger
);

-- Избранные товары
CREATE TABLE IF NOT EXISTS bazaar.favorite
(
    id         UUID PRIMARY KEY,
    user_id    UUID REFERENCES bazaar."user" (id) ON DELETE CASCADE,
    product_id UUID REFERENCES bazaar.product (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (user_id, product_id)
);

-- Картинки товаров
CREATE TABLE IF NOT EXISTS bazaar.product_image
(
    id         UUID PRIMARY KEY,
    product_id UUID REFERENCES bazaar.product (id) ON DELETE CASCADE,
    image_url  TEXT NOT NULL,
    num        INT CHECK (num >= 0)
);

-- Скидки
CREATE TABLE IF NOT EXISTS bazaar.discount
(
    id               UUID PRIMARY KEY,
    start_date       TIMESTAMPTZ NOT NULL,
    end_date         TIMESTAMPTZ NOT NULL,
    product_id       UUID REFERENCES bazaar.product (id) ON DELETE CASCADE,
    discounted_price NUMERIC(12, 2) CHECK (discounted_price >= 0),
    updated_at       TIMESTAMPTZ DEFAULT now()
);

-- Категории
CREATE TABLE IF NOT EXISTS bazaar.category
(
    id   UUID PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Привязка товаров к категориям
CREATE TABLE IF NOT EXISTS bazaar.product_category
(
    id          UUID PRIMARY KEY,
    product_id  UUID REFERENCES bazaar.product (id) ON DELETE CASCADE,
    category_id UUID REFERENCES bazaar.category (id) ON DELETE CASCADE,
    UNIQUE (product_id, category_id)
);

-- Заказы
CREATE TABLE IF NOT EXISTS bazaar."order"
(
    id                   UUID PRIMARY KEY,
    user_id              UUID REFERENCES bazaar."user" (id) ON DELETE CASCADE,
    status               bazaar.order_status                              NOT NULL,
    total_price          NUMERIC(12, 2) CHECK (total_price >= 0)          NOT NULL,
    total_price_discount NUMERIC(12, 2) CHECK (total_price_discount >= 0) NOT NULL,
    address_id           UUID                                             REFERENCES bazaar.address (id) ON DELETE SET NULL,
    created_at           TIMESTAMPTZ DEFAULT now(),
    updated_at           TIMESTAMPTZ DEFAULT now()
);

-- Элементы заказа
CREATE TABLE IF NOT EXISTS bazaar.order_item
(
    id         UUID PRIMARY KEY,
    order_id   UUID REFERENCES bazaar."order" (id) ON DELETE CASCADE,
    product_id UUID REFERENCES bazaar.product (id) ON DELETE CASCADE,
    price      NUMERIC(12, 2) CHECK (price >= 0) NOT NULL,
    quantity   INT CHECK (quantity > 0)          NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (order_id, product_id)
);

-- Корзина
CREATE TABLE IF NOT EXISTS bazaar.basket
(
    id                   UUID PRIMARY KEY,
    user_id              UUID REFERENCES bazaar."user" (id) ON DELETE CASCADE UNIQUE,
    total_price          NUMERIC(12, 2) CHECK (total_price >= 0)          NOT NULL,
    total_price_discount NUMERIC(12, 2) CHECK (total_price_discount >= 0) NOT NULL
);

-- Элементы корзины
CREATE TABLE IF NOT EXISTS bazaar.basket_item
(
    id         UUID PRIMARY KEY,
    basket_id  UUID REFERENCES bazaar.basket (id) ON DELETE CASCADE,
    product_id UUID REFERENCES bazaar.product (id) ON DELETE CASCADE,
    quantity   INT CHECK (quantity > 0) NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE (basket_id, product_id)
);

-- Отзывы
CREATE TABLE IF NOT EXISTS bazaar.review
(
    id         UUID PRIMARY KEY,
    user_id    UUID REFERENCES bazaar."user" (id) ON DELETE CASCADE,
    product_id UUID REFERENCES bazaar.product (id) ON DELETE CASCADE,
    rating     INT CHECK (rating BETWEEN 1 AND 5) NOT NULL,
    comment    TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Промокоды
CREATE TABLE IF NOT EXISTS bazaar.promo_code
(
    id                UUID PRIMARY KEY,
    user_id           UUID REFERENCES bazaar."user" (id) ON DELETE CASCADE,
    category_id       UUID REFERENCES bazaar.category (id) ON DELETE CASCADE,
    seller_id         UUID REFERENCES bazaar."user" (id) ON DELETE CASCADE,
    code              TEXT UNIQUE NOT NULL,
    relative_discount INT CHECK (relative_discount BETWEEN 0 AND 1),
    absolute_discount NUMERIC(12, 2) CHECK (absolute_discount >= 0),
    start_date        TIMESTAMPTZ NOT NULL,
    end_date          TIMESTAMPTZ NOT NULL,
    updated_at        TIMESTAMPTZ DEFAULT now()
);