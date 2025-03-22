## Таблицы

### `user`
Хранит информацию о пользователях платформы.
- `id` — уникальный идентификатор
- `email` — почта пользователя
- `phone_number` — номер телефона
- `password_hash` — хеш пароля
- `name` — имя пользователя
- `surname` — фамилия пользователя
- `image_url` — ссылка на изображение профиля
- `role` — роль пользователя (`seller`, `buyer`, `admin`)
- `version` — версия сессии пользователя

### `user_role`
Определяет возможные роли пользователей.
- `seller` — продавец
- `buyer` — покупатель
- `admin` — администратор

### `product`
Содержит информацию о товарах, продаваемых на платформе.
- `id` — уникальный идентификатор
- `seller_id` — идентификатор продавца
- `name` — название товара
- `image_url` — изображение товара
- `description` — описание товара
- `price` — цена
- `category_id` — категория товара
- `quantity` — количество товара в наличии

### `discount`
Содержит информацию о скидках на товары.
- `id` — уникальный идентификатор
- `end_date` — дата окончания скидки
- `product_id` — товар, на который действует скидка
- `discounted_price` — цена со скидкой

### `category`
Хранит категории товаров.
- `id` — уникальный идентификатор
- `name` — название категории
- `image_url` — изображение категории

### `order`
Хранит информацию о заказах.
- `id` — уникальный идентификатор
- `user_id` — пользователь, сделавший заказ
- `status` — статус заказа (`pending`, `paid`, `shipped`, `delivered`, `canceled`)
- `total_price` — общая сумма заказа
- `address_id` — адрес доставки
- `created_at` — дата создания заказа

### `order_status`
Определяет возможные статусы заказов.
- `pending` — в ожидании
- `paid` — оплачено
- `shipped` — отправлено
- `delivered` — доставлено
- `canceled` — отменено

### `order_item`
Связывает заказы с товарами.
- `id` — уникальный идентификатор
- `order_id` — идентификатор заказа
- `product_id` — идентификатор товара
- `quantity` — количество товара в заказе

### `basket`
Корзина пользователя.
- `id` — уникальный идентификатор
- `user_id` — идентификатор владельца корзины

### `basket_item`
Связывает товары с корзиной пользователя.
- `id` — уникальный идентификатор
- `basket_id` — идентификатор корзины
- `product_id` — идентификатор товара
- `quantity` — количество товара в корзине

### `review`
Хранит отзывы пользователей о товарах.
- `id` — уникальный идентификатор
- `user_id` — идентификатор пользователя, оставившего отзыв
- `product_id` — идентификатор товара
- `rating` — оценка товара
- `comment` — комментарий пользователя
- `created_at` — дата отзыва

### `address`
Хранит адреса пользователей.
- `id` — уникальный идентификатор
- `user_id` — идентификатор владельца адреса
- `city` — город
- `street` — улица
- `zip_code` — почтовый индекс

---

## ER-диаграмма базы данных

```mermaid
erDiagram

    user {
        id PK
        email
        phone_number
        password_hash
        name
        surname
        image_url
        role FK
        version
    }

    user_role {
        seller
        buyer
        admin
    }

    product {
        id PK
        seller_id FK
        name
        image_url
        description
        price
        category_id FK
        quantity
    }

    discount {
        id PK
        end_date
        product_id FK
        discounted_price
    }

    category {
        id PK
        name
        image_url
    }

    order {
        id PK
        user_id FK
        status
        total_price
        address_id FK
        created_at
    }

    order_status {
        pending
        paid
        shipped
        delivered
        canceled
    }

    order_item {
        id PK
        order_id FK
        product_id FK
        quantity
    }

    basket {
        id PK
        user_id FK
    }

    basket_item {
        id PK
        basket_id FK
        product_id FK
        quantity
    }

    review {
        id PK
        user_id FK
        product_id FK
        rating
        comment
        created_at
    }

    address {
        id PK
        user_id FK
        city
        street
        zip_code
    }

    user ||--o{ order : "создает"
    user ||--o{ address : "имеет"
    user ||--o{ product : "продает"
    user ||--o{ review : "оставляет"
    user ||--o{ basket : "имеет"

    product ||--o{ order_item : "включается в"
    product ||--o{ review : "получает"
    product }o--|| category : "принадлежит"
    product ||--o{ discount : "имеет"

    order ||--o{ order_item : "содержит"
    order }o--|| address : "доставляется в"

    basket ||--o{ basket_item : "содержит"
    basket_item }o--|| product : "относится к"
