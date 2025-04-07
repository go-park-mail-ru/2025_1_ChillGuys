CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_version_updated_at
    BEFORE UPDATE ON bazaar.user_version
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_product_updated_at
    BEFORE UPDATE ON bazaar.product
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_discount_updated_at
    BEFORE UPDATE ON bazaar.discount
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_promo_code_updated_at
    BEFORE UPDATE ON bazaar.promo_code
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_order_updated_at
    BEFORE UPDATE ON bazaar."order"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_order_item_updated_at
    BEFORE UPDATE ON bazaar.order_item
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_basket_item_updated_at
    BEFORE UPDATE ON bazaar.basket_item
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_review_updated_at
    BEFORE UPDATE ON bazaar.review
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_pickup_point_updated_at
    BEFORE UPDATE ON bazaar.pickup_point
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_address_updated_at
    BEFORE UPDATE ON bazaar.address
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER update_user_balance_updated_at
    BEFORE UPDATE ON bazaar.user_balance
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at();