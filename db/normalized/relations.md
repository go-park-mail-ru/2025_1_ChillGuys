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
- `role` — роль пользователя 
- `version` — версия сессии пользователя

### `product`
Содержит информацию о товарах.
- `id` — уникальный идентификатор
- `seller_id` — идентификатор продавца
- `name` — название товара
- `image_url` — изображение товара
- `description` — описание товара
- `price` — цена
- `category_id` — категория товара
- `quantity` — количество товара в наличии

### `product_image`
Хранит изображения товаров.
- `id` — уникальный идентификатор
- `product_id` — идентификатор продавца
- `image_url` — ссылка на изображение товара

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

### `order`
Хранит информацию о заказах.
- `id` — уникальный идентификатор
- `user_id` — пользователь, сделавший заказ
- `status` — статус
- `total_price` — общая сумма заказа
- `address_id` — адрес доставки
- `created_at` — дата создания заказа

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
        uuid id PK
        string email
        string phone_number
        string password_hash
        string name
        string surname
        string image_url
        string role
        int version
    }

    product {
        uuid id PK
        uuid seller_id FK
        string name
        string preview_image_url
        string description
        float price
        uuid category_id FK
        int quantity
    }
    
    product_image {
        uuid id PK
        uuid product_id FK
        string image_url
    }

    discount {
        uuid id PK
        datetime end_date
        uuid product_id FK
        float discounted_price
    }

    category {
        uuid id PK
        string name
    }

    order {
        uuid id PK
        uuid user_id FK
        string status
        float total_price
        uuid address_id FK
        datetime created_at
    }
    
    order_item {
        uuid id PK
        uuid order_id FK
        uuid product_id FK
        int quantity
    }

    basket {
        uuid id PK
        uuid user_id FK
    }

    basket_item {
        uuid id PK
        uuid basket_id FK
        uuid product_id FK
        int quantity
    }

    review {
        uuid id PK
        uuid user_id FK
        uuid product_id FK
        int rating
        string comment
        datetime created_at
    }

    address {
        uuid id PK
        uuid user_id FK
        string city
        string street
        string zip_code
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
    product ||--o{ product_image : "имеет"

    order ||--o{ order_item : "содержит"
    order }o--|| address : "доставляется в"

    basket ||--o{ basket_item : "содержит"
    basket_item }o--|| product : "относится к"
```