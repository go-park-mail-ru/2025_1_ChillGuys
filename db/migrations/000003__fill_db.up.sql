-- Вставка адресов
INSERT INTO bazaar.address (id, city, street, house, zip_code)
VALUES ('550e8400-e29b-41d4-a716-446655440201', 'Москва', 'Ленина', '1', '101000'),
       ('550e8400-e29b-41d4-a716-446655440202', 'Санкт-Петербург', 'Невский', '10', '190000'),
       ('550e8400-e29b-41d4-a716-446655440203', 'Москва', 'Ленина', '1', '101000');


-- Вставка ПВЗ
INSERT INTO bazaar.pickup_point (id, address_id)
VALUES ('550e8400-e29b-41d4-a716-446655440040', '550e8400-e29b-41d4-a716-446655440201'),
       ('550e8400-e29b-41d4-a716-446655440041', '550e8400-e29b-41d4-a716-446655440202'),
       ('550e8400-e29b-41d4-a716-446655440042', '550e8400-e29b-41d4-a716-446655440203');

-- Вставка пользователей
INSERT INTO bazaar."user" (id, email, phone_number, password_hash, name)
VALUES ('e29b41d4-a716-4466-5544-000000000001', 'user1@example.com', '+1234567890', 'hash1', 'Иван'),
       ('e29b41d4-a716-4466-5544-000000000002', 'user2@example.com', '+1234567891', 'hash2', 'Петр'),
       ('e29b41d4-a716-4466-5544-000000000003', 'seller1@example.com', '+1234567892', 'hash3', 'Алексей'),
       ('e29b41d4-a716-4466-5544-000000000004', 'seller2@example.com', '+1234567893', 'hash4', 'Мария'),
       ('e29b41d4-a716-4466-5544-000000000005', 'seller3@example.com', '+1234567894', 'hash5', 'Ольга'),
       ('e29b41d4-a716-4466-5544-000000000006', 'seller4@example.com', '+1234567895', 'hash6', 'Дмитрий'),
       ('e29b41d4-a716-4466-5544-000000000007', 'seller5@example.com', '+1234567896', 'hash7', 'Анна'),
       ('e29b41d4-a716-4466-5544-000000000008', 'seller6@example.com', '+1234567897', 'hash8', 'Сергей'),
       ('e29b41d4-a716-4466-5544-000000000009', 'seller7@example.com', '+1234567898', 'hash9', 'Елена'),
       ('e29b41d4-a716-4466-5544-000000000010', 'seller8@example.com', '+1234567899', 'hash10', 'Андрей');

-- Вставка ролей
INSERT INTO bazaar.role (id, name)
VALUES ('550e8400-e29b-41d4-a716-446655440001', 'admin'),
       ('550e8400-e29b-41d4-a716-446655440002', 'seller'),
       ('550e8400-e29b-41d4-a716-446655440003', 'customer');

-- Вставка связи пользователей с ролями
INSERT INTO bazaar.user_role (id, user_id, role_id)
VALUES ('550e8400-e29b-41d4-a716-446655440011', 'e29b41d4-a716-4466-5544-000000000001',
        '550e8400-e29b-41d4-a716-446655440003'),
       ('550e8400-e29b-41d4-a716-446655440012', 'e29b41d4-a716-4466-5544-000000000002',
        '550e8400-e29b-41d4-a716-446655440003'),
       ('550e8400-e29b-41d4-a716-446655440013', 'e29b41d4-a716-4466-5544-000000000003',
        '550e8400-e29b-41d4-a716-446655440002'),
       ('550e8400-e29b-41d4-a716-446655440014', 'e29b41d4-a716-4466-5544-000000000004',
        '550e8400-e29b-41d4-a716-446655440002'),
       ('550e8400-e29b-41d4-a716-446655440015', 'e29b41d4-a716-4466-5544-000000000005',
        '550e8400-e29b-41d4-a716-446655440002'),
       ('550e8400-e29b-41d4-a716-446655440016', 'e29b41d4-a716-4466-5544-000000000006',
        '550e8400-e29b-41d4-a716-446655440002'),
       ('550e8400-e29b-41d4-a716-446655440017', 'e29b41d4-a716-4466-5544-000000000007',
        '550e8400-e29b-41d4-a716-446655440002'),
       ('550e8400-e29b-41d4-a716-446655440018', 'e29b41d4-a716-4466-5544-000000000008',
        '550e8400-e29b-41d4-a716-446655440002'),
       ('550e8400-e29b-41d4-a716-446655440019', 'e29b41d4-a716-4466-5544-000000000009',
        '550e8400-e29b-41d4-a716-446655440002'),
       ('550e8400-e29b-41d4-a716-446655440020', 'e29b41d4-a716-4466-5544-000000000010',
        '550e8400-e29b-41d4-a716-446655440002');

-- Вставка категорий
INSERT INTO bazaar.category (id, name)
VALUES ('550e8400-e29b-41d4-a716-446655440010', 'Электроника'),
       ('550e8400-e29b-41d4-a716-446655440011', 'Бытовая техника'),
       ('550e8400-e29b-41d4-a716-446655440012', 'Компьютеры и ноутбуки'),
       ('550e8400-e29b-41d4-a716-446655440013', 'Смартфоны и гаджеты'),
       ('550e8400-e29b-41d4-a716-446655440014', 'Аудиотехника'),
       ('550e8400-e29b-41d4-a716-446655440015', 'Игровые консоли');

-- Вставка всех 15 товаров
INSERT INTO bazaar.product (id, seller_id, name, description, status, price, quantity, rating, reviews_count)
VALUES ('550e8400-e29b-41d4-a716-446655440001', 'e29b41d4-a716-4466-5544-000000000003', 'Смартфон Xiaomi Redmi Note 10',
        'Смартфон с AMOLED-дисплеем и камерой 48 Мп', 'approved', 19999.00, 50, 4, 120),
       ('550e8400-e29b-41d4-a716-446655440002', 'e29b41d4-a716-4466-5544-000000000004', 'Ноутбук ASUS VivoBook 15',
        'Ноутбук с процессором Intel Core i5 и SSD на 512 ГБ', 'approved', 54999.00, 30, 4, 80),
       ('550e8400-e29b-41d4-a716-446655440003', 'e29b41d4-a716-4466-5544-000000000005', 'Наушники Sony WH-1000XM4',
        'Беспроводные наушники с шумоподавлением', 'approved', 29999.00, 25, 4, 200),
       ('550e8400-e29b-41d4-a716-446655440004', 'e29b41d4-a716-4466-5544-000000000006',
        'Фитнес-браслет Xiaomi Mi Band 6', 'Фитнес-браслет с AMOLED-дисплеем и мониторингом сна', 'approved', 3999.00,
        100, 4, 300),
       ('550e8400-e29b-41d4-a716-446655440005', 'e29b41d4-a716-4466-5544-000000000007', 'Пылесос Dyson V11',
        'Беспроводной пылесос с мощным всасыванием', 'approved', 59999.00, 15, 4, 90),
       ('550e8400-e29b-41d4-a716-446655440006', 'e29b41d4-a716-4466-5544-000000000008', 'Кофемашина DeLonghi Magnifica',
        'Автоматическая кофемашина для приготовления эспрессо', 'approved', 79999.00, 10, 4, 70),
       ('550e8400-e29b-41d4-a716-446655440007', 'e29b41d4-a716-4466-5544-000000000009',
        'Электросамокат Xiaomi Mi Scooter 3', 'Электросамокат с запасом хода 30 км', 'approved', 29999.00, 40, 4, 150),
       ('550e8400-e29b-41d4-a716-446655440008', 'e29b41d4-a716-4466-5544-000000000010',
        'Умная колонка Яндекс.Станция Мини', 'Умная колонка с голосовым помощником Алисой', 'approved', 7999.00, 60, 4,
        250),
       ('550e8400-e29b-41d4-a716-446655440009', 'e29b41d4-a716-4466-5544-000000000003', 'Монитор Samsung Odyssey G5',
        'Игровой монитор с разрешением 1440p и частотой 144 Гц', 'approved', 34999.00, 20, 4, 100),
       ('550e8400-e29b-41d4-a716-446655440010', 'e29b41d4-a716-4466-5544-000000000004', 'Электрочайник Bosch TWK 3A011',
        'Электрочайник с мощностью 2400 Вт', 'approved', 1999.00, 50, 4, 180),
       ('550e8400-e29b-41d4-a716-446655440011', 'e29b41d4-a716-4466-5544-000000000005',
        'Робот-пылесос iRobot Roomba 981', 'Робот-пылесос с навигацией по карте помещения', 'approved', 69999.00, 12, 4,
        60),
       ('550e8400-e29b-41d4-a716-446655440012', 'e29b41d4-a716-4466-5544-000000000006', 'Фен Dyson Supersonic',
        'Фен с технологией защиты волос от перегрева', 'approved', 49999.00, 18, 4, 130),
       ('550e8400-e29b-41d4-a716-446655440013', 'e29b41d4-a716-4466-5544-000000000007',
        'Микроволновая печь LG MS-2042DB', 'Микроволновка с объемом 20 литров', 'approved', 8999.00, 35, 4, 110),
       ('550e8400-e29b-41d4-a716-446655440014', 'e29b41d4-a716-4466-5544-000000000008', 'Игровая консоль PlayStation 5',
        'Игровая консоль нового поколения', 'approved', 79999.00, 5, 4, 300),
       ('550e8400-e29b-41d4-a716-446655440015', 'e29b41d4-a716-4466-5544-000000000009',
        'Электронная книга PocketBook 740', 'Электронная книга с экраном E Ink Carta', 'approved', 19999.00, 25, 4, 90);

-- Вставка связи товаров с категориями
INSERT INTO bazaar.product_category (id, product_id, category_id)
VALUES ('550e8400-e29b-41d4-a716-446655440101', '550e8400-e29b-41d4-a716-446655440001',
        '550e8400-e29b-41d4-a716-446655440013'),
       ('550e8400-e29b-41d4-a716-446655440102', '550e8400-e29b-41d4-a716-446655440002',
        '550e8400-e29b-41d4-a716-446655440012'),
       ('550e8400-e29b-41d4-a716-446655440103', '550e8400-e29b-41d4-a716-446655440003',
        '550e8400-e29b-41d4-a716-446655440014'),
       ('550e8400-e29b-41d4-a716-446655440104', '550e8400-e29b-41d4-a716-446655440004',
        '550e8400-e29b-41d4-a716-446655440013'),
       ('550e8400-e29b-41d4-a716-446655440105', '550e8400-e29b-41d4-a716-446655440005',
        '550e8400-e29b-41d4-a716-446655440011'),
       ('550e8400-e29b-41d4-a716-446655440106', '550e8400-e29b-41d4-a716-446655440006',
        '550e8400-e29b-41d4-a716-446655440011'),
       ('550e8400-e29b-41d4-a716-446655440107', '550e8400-e29b-41d4-a716-446655440007',
        '550e8400-e29b-41d4-a716-446655440010'),
       ('550e8400-e29b-41d4-a716-446655440108', '550e8400-e29b-41d4-a716-446655440008',
        '550e8400-e29b-41d4-a716-446655440010'),
       ('550e8400-e29b-41d4-a716-446655440109', '550e8400-e29b-41d4-a716-446655440009',
        '550e8400-e29b-41d4-a716-446655440012'),
       ('550e8400-e29b-41d4-a716-446655440110', '550e8400-e29b-41d4-a716-446655440010',
        '550e8400-e29b-41d4-a716-446655440011'),
       ('550e8400-e29b-41d4-a716-446655440111', '550e8400-e29b-41d4-a716-446655440011',
        '550e8400-e29b-41d4-a716-446655440011'),
       ('550e8400-e29b-41d4-a716-446655440112', '550e8400-e29b-41d4-a716-446655440012',
        '550e8400-e29b-41d4-a716-446655440011'),
       ('550e8400-e29b-41d4-a716-446655440113', '550e8400-e29b-41d4-a716-446655440013',
        '550e8400-e29b-41d4-a716-446655440011'),
       ('550e8400-e29b-41d4-a716-446655440114', '550e8400-e29b-41d4-a716-446655440014',
        '550e8400-e29b-41d4-a716-446655440015'),
       ('550e8400-e29b-41d4-a716-446655440115', '550e8400-e29b-41d4-a716-446655440015',
        '550e8400-e29b-41d4-a716-446655440010');

-- Вставка балансов пользователей
INSERT INTO bazaar.user_balance (id, user_id, balance)
VALUES ('550e8400-e29b-41d4-a716-446655440301', 'e29b41d4-a716-4466-5544-000000000001', 100000.00),
       ('550e8400-e29b-41d4-a716-446655440302', 'e29b41d4-a716-4466-5544-000000000002', 50000.00);

-- Вставка корзин
INSERT INTO bazaar.basket (id, user_id, total_price, total_price_discount)
VALUES ('550e8400-e29b-41d4-a716-446655440401', 'e29b41d4-a716-4466-5544-000000000001', 0.00, 0.00),
       ('550e8400-e29b-41d4-a716-446655440402', 'e29b41d4-a716-4466-5544-000000000002', 0.00, 0.00);

-- Вставка элементов корзины
INSERT INTO bazaar.basket_item (id, basket_id, product_id, quantity)
VALUES ('550e8400-e29b-41d4-a716-446655440501', '550e8400-e29b-41d4-a716-446655440401',
        '550e8400-e29b-41d4-a716-446655440001', 1),
       ('550e8400-e29b-41d4-a716-446655440502', '550e8400-e29b-41d4-a716-446655440401',
        '550e8400-e29b-41d4-a716-446655440003', 1),
       ('550e8400-e29b-41d4-a716-446655440503', '550e8400-e29b-41d4-a716-446655440402',
        '550e8400-e29b-41d4-a716-446655440005', 1);

-- Вставка заказов
INSERT INTO bazaar."order" (id, user_id, status, total_price, total_price_discount, address_id)
VALUES ('550e8400-e29b-41d4-a716-446655440601', 'e29b41d4-a716-4466-5544-000000000001', 'placed', 49998.00, 49998.00,
        '550e8400-e29b-41d4-a716-446655440201'),
       ('550e8400-e29b-41d4-a716-446655440602', 'e29b41d4-a716-4466-5544-000000000002', 'delivered', 59999.00, 59999.00,
        '550e8400-e29b-41d4-a716-446655440202');

-- Вставка элементов заказа
INSERT INTO bazaar.order_item (id, order_id, product_id, price, quantity)
VALUES ('550e8400-e29b-41d4-a716-446655440701', '550e8400-e29b-41d4-a716-446655440601',
        '550e8400-e29b-41d4-a716-446655440001', 1499.99, 1),
       ('550e8400-e29b-41d4-a716-446655440702', '550e8400-e29b-41d4-a716-446655440601',
        '550e8400-e29b-41d4-a716-446655440003', 2599.00, 1),
       ('550e8400-e29b-41d4-a716-446655440703', '550e8400-e29b-41d4-a716-446655440602',
        '550e8400-e29b-41d4-a716-446655440005', 799.50, 1);

-- Вставка избранных товаров
INSERT INTO bazaar.favorite (id, user_id, product_id)
VALUES ('550e8400-e29b-41d4-a716-446655440801', 'e29b41d4-a716-4466-5544-000000000001',
        '550e8400-e29b-41d4-a716-446655440002'),
       ('550e8400-e29b-41d4-a716-446655440802', 'e29b41d4-a716-4466-5544-000000000001',
        '550e8400-e29b-41d4-a716-446655440004'),
       ('550e8400-e29b-41d4-a716-446655440803', 'e29b41d4-a716-4466-5544-000000000002',
        '550e8400-e29b-41d4-a716-446655440001'),
       ('550e8400-e29b-41d4-a716-446655440804', 'e29b41d4-a716-4466-5544-000000000002',
        '550e8400-e29b-41d4-a716-446655440006');

-- Вставка отзывов
INSERT INTO bazaar.review (id, user_id, product_id, rating, comment)
VALUES ('550e8400-e29b-41d4-a716-446655440901', 'e29b41d4-a716-4466-5544-000000000001',
        '550e8400-e29b-41d4-a716-446655440001', 5, 'Отличный смартфон за свои деньги!'),
       ('550e8400-e29b-41d4-a716-446655440902', 'e29b41d4-a716-4466-5544-000000000001',
        '550e8400-e29b-41d4-a716-446655440003', 4, 'Хорошие наушники, но дороговаты'),
       ('550e8400-e29b-41d4-a716-446655440903', 'e29b41d4-a716-4466-5544-000000000002',
        '550e8400-e29b-41d4-a716-446655440005', 5, 'Лучший пылесос, который у меня был!'),
       ('550e8400-e29b-41d4-a716-446655440904', 'e29b41d4-a716-4466-5544-000000000002',
        '550e8400-e29b-41d4-a716-446655440008', 4, 'Хороший звук, удобное управление');

-- Вставка скидок
INSERT INTO bazaar.discount (id, start_date, end_date, product_id, discounted_price)
VALUES ('550e8400-e29b-41d4-a716-446655441001', NOW(), NOW() + INTERVAL '7 DAY', '550e8400-e29b-41d4-a716-446655440001',
        17999.00),
       ('550e8400-e29b-41d4-a716-446655441002', NOW(), NOW() + INTERVAL '14 DAY',
        '550e8400-e29b-41d4-a716-446655440003', 26999.00),
       ('550e8400-e29b-41d4-a716-446655441003', NOW(), NOW() + INTERVAL '10 DAY',
        '550e8400-e29b-41d4-a716-446655440005', 54999.00);

-- Вставка промокодов
INSERT INTO bazaar.promo_code (id, code, relative_discount, absolute_discount, start_date, end_date)
VALUES ('550e8400-e29b-41d4-a716-446655441101', 'SUMMER10', 0.1, NULL, NOW(), NOW() + INTERVAL '30 DAY'),
       ('550e8400-e29b-41d4-a716-446655441102', 'TECH500', NULL, 500.00, NOW(), NOW() + INTERVAL '15 DAY'),
       ('550e8400-e29b-41d4-a716-446655441103', 'XIAOMI20', 0.2, NULL, NOW(), NOW() + INTERVAL '7 DAY');

-- Вставка версий пользователей
INSERT INTO bazaar.user_version (id, user_id)
VALUES ('550e8400-e29b-41d4-a716-446655441201', 'e29b41d4-a716-4466-5544-000000000001'),
       ('550e8400-e29b-41d4-a716-446655441202', 'e29b41d4-a716-4466-5544-000000000002');

-- Вставка изображений товаров
INSERT INTO bazaar.product_image (id, product_id, image_url, num)
VALUES ('550e8400-e29b-41d4-a716-446655441301', '550e8400-e29b-41d4-a716-446655440001', 'xiaomi_redmi_note_10_1.jpg',
        0),
       ('550e8400-e29b-41d4-a716-446655441302', '550e8400-e29b-41d4-a716-446655440001', 'xiaomi_redmi_note_10_2.jpg',
        1),
       ('550e8400-e29b-41d4-a716-446655441303', '550e8400-e29b-41d4-a716-446655440002', 'asus_vivobook_15_1.jpg', 0),
       ('550e8400-e29b-41d4-a716-446655441304', '550e8400-e29b-41d4-a716-446655440003', 'sony_wh1000xm4_1.jpg', 0),
       ('550e8400-e29b-41d4-a716-446655441305', '550e8400-e29b-41d4-a716-446655440004', 'xiaomi_mi_band_6_1.jpg', 0),
       ('550e8400-e29b-41d4-a716-446655441306', '550e8400-e29b-41d4-a716-446655440005', 'dyson_v11_1.jpg', 0),
       ('550e8400-e29b-41d4-a716-446655441307', '550e8400-e29b-41d4-a716-446655440006', 'delonghi_magnifica_1.jpg', 0),
       ('550e8400-e29b-41d4-a716-446655441308', '550e8400-e29b-41d4-a716-446655440007', 'xiaomi_mi_scooter_3_1.jpg', 0),
       ('550e8400-e29b-41d4-a716-446655441309', '550e8400-e29b-41d4-a716-446655440008', 'yandex_station_mini_1.jpg', 0),
       ('550e8400-e29b-41d4-a716-446655441310', '550e8400-e29b-41d4-a716-446655440009', 'samsung_odyssey_g5_1.jpg', 0),
       ('550e8400-e29b-41d4-a716-446655441311', '550e8400-e29b-41d4-a716-446655440010', 'bosch_twk_3a011_1.jpg', 0),
       ('550e8400-e29b-41d4-a716-446655441312', '550e8400-e29b-41d4-a716-446655440011', 'irobot_roomba_981_1.jpg', 0),
       ('550e8400-e29b-41d4-a716-446655441313', '550e8400-e29b-41d4-a716-446655440012', 'dyson_supersonic_1.jpg', 0),
       ('550e8400-e29b-41d4-a716-446655441314', '550e8400-e29b-41d4-a716-446655440013', 'lg_ms2042db_1.jpg', 0),
       ('550e8400-e29b-41d4-a716-446655441315', '550e8400-e29b-41d4-a716-446655440014', 'playstation5_1.jpg', 0),
       ('550e8400-e29b-41d4-a716-446655441316', '550e8400-e29b-41d4-a716-446655440015', 'pocketbook_740_1.jpg', 0);