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

### `user_role`
Хранит информацию о ролях пользователей.
- `id` — уникальный идентификатор
- `user_id` — идентификатор пользователя
- `role` — роль пользователя

### `user_balance`
Хранит баланс пользователей.
- `id` — уникальный идентификатор
- `user_id` — идентификатор пользователя
- `balance` — текущий баланс
- `updated_at` — время последнего обновления

### `user_version`
Хранит версию сессии пользователя.
- `id` — уникальный идентификатор
- `user_id` — идентификатор пользователя
- `version` — версия сессии
- `updated_at` — время последнего обновления

### `product`
Содержит информацию о товарах.
- `id` — уникальный идентификатор
- `seller_id` — идентификатор продавца
- `name` — название товара
- `preview_image_url` — изображение товара
- `description` — описание товара
- `price` — цена
- `category_id` — категория товара
- `quantity` — количество товара в наличии
- `updated_at` — время последнего обновления

### `product_image`
Хранит изображения товаров.
- `id` — уникальный идентификатор
- `product_id` — идентификатор товара
- `image_url` — ссылка на изображение товара
- `num` — порядковый номер изображения

### `discount`
Содержит информацию о скидках на товары.
- `id` — уникальный идентификатор
- `start_date` — дата начала скидки
- `end_date` — дата окончания скидки
- `product_id` — товар, на который действует скидка
- `discounted_price` — цена со скидкой
- `updated_at` — время последнего обновления

### `promo_code`
Хранит информацию о промокодах.
- `id` — уникальный идентификатор
- `user_id` — идентификатор пользователя, который может воспользоваться промокодом
- `category_id` — идентификатор категории (если промокод на товары конкретной категории)
- `seller_id` — идентификатор продавца (если промокод на товары конкретного продавца)
- `code` — значение промокода
- `relative_discount` — скидка в процентах
- `absolute_discount` — скидка в абсолютном значении
- `start_date` — дата начала действия
- `end_date` — дата окончания действия
- `updated_at` — время последнего обновления

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
- `updated_at` — время последнего обновления

### `order_item`
Связывает заказы с товарами.
- `id` — уникальный идентификатор
- `order_id` — идентификатор заказа
- `product_id` — идентификатор товара
- `quantity` — количество товара в заказе
- `updated_at` — время последнего обновления

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
- `updated_at` — время последнего обновления

### `review`
Хранит отзывы пользователей о товарах.
- `id` — уникальный идентификатор
- `user_id` — идентификатор пользователя, оставившего отзыв
- `product_id` — идентификатор товара
- `rating` — оценка товара
- `comment` — комментарий пользователя
- `created_at` — дата отзыва
- `updated_at` — время последнего обновления

### `address`
Хранит адреса пользователей.
- `id` — уникальный идентификатор
- `user_id` — идентификатор владельца адреса
- `city` — город
- `street` — улица
- `house` — номер дома
- `apartment` — номер квартиры
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
    }

    user_role {
        uuid id PK
        uuid user_id FK
        string role
    }

    user_balance {
        uuid id PK
        uuid user_id FK
        float balance
        datetime updated_at
    }

    user_version {
        uuid id PK
        uuid user_id FK
        int version
        datetime updated_at
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
        datetime updated_at
    }
    
    product_image {
        uuid id PK
        uuid product_id FK
        string image_url
        int num
    }

    discount {
        uuid id PK
        datetime start_date
        datetime end_date
        uuid product_id FK
        float discounted_price
        datetime updated_at
    }

    promo_code {
        uuid id PK
        uuid user_id FK
        uuid category_id FK 
        uuid seller_id FK 
        string code
        float relative_discount 
        float absolute_discount 
        datetime start_date
        datetime end_date
        datetime updated_at
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
        datetime updated_at
    }
    
    order_item {
        uuid id PK
        uuid order_id FK
        uuid product_id FK
        int quantity
        datetime updated_at
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
        datetime updated_at
    }

    review {
        uuid id PK
        uuid user_id FK
        uuid product_id FK
        int rating
        string comment
        datetime created_at
        datetime updated_at
    }

    address {
        uuid id PK
        uuid user_id FK
        string city
        string street
        string house
        string apartment
        string zip_code
    }

    user ||--o{ user_role : "имеет роль"
    user ||--o{ user_balance : "имеет баланс"
    user ||--o{ user_version : "имеет версию"
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

    promo_code ||--|| user : "принадлежит"
    promo_code ||--|| category : "может применяться к категории"
    promo_code ||--|| user : "может применяться к продавцу"
```