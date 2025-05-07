CREATE OR REPLACE FUNCTION update_updated_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_version_updated_at
    BEFORE UPDATE
    ON bazaar.user_version
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_product_updated_at
    BEFORE UPDATE
    ON bazaar.product
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_discount_updated_at
    BEFORE UPDATE
    ON bazaar.discount
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_promo_code_updated_at
    BEFORE UPDATE
    ON bazaar.promo_code
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_order_updated_at
    BEFORE UPDATE
    ON bazaar."order"
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_order_item_updated_at
    BEFORE UPDATE
    ON bazaar.order_item
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_basket_item_updated_at
    BEFORE UPDATE
    ON bazaar.basket_item
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_review_updated_at
    BEFORE UPDATE
    ON bazaar.review
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_pickup_point_updated_at
    BEFORE UPDATE
    ON bazaar.pickup_point
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_address_updated_at
    BEFORE UPDATE
    ON bazaar.address
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_user_balance_updated_at
    BEFORE UPDATE
    ON bazaar.user_balance
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

CREATE OR REPLACE FUNCTION update_product_review_stats()
RETURNS TRIGGER AS $$
BEGIN
    -- Обновляем статистику продукта при любых изменениях в отзывах
    UPDATE bazaar.product p
    SET 
        reviews_count = subquery.review_count,
        rating = subquery.avg_rating
    FROM (
        SELECT 
            COUNT(*) as review_count,
            COALESCE(AVG(rating), 0) as avg_rating
        FROM bazaar.review
        WHERE product_id = COALESCE(NEW.product_id, OLD.product_id)
    ) subquery
    WHERE p.id = COALESCE(NEW.product_id, OLD.product_id);
    
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Триггер для вставки нового отзыва
CREATE TRIGGER after_review_insert
AFTER INSERT ON bazaar.review
FOR EACH ROW
EXECUTE FUNCTION update_product_review_stats();
