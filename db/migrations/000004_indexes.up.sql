CREATE UNIQUE INDEX IF NOT EXISTS idx_address_unique
    ON bazaar.address (city, street, house, apartment, zip_code);

CREATE UNIQUE INDEX IF NOT EXISTS user_address_unique
    ON bazaar.user_address (user_id, address_id);
