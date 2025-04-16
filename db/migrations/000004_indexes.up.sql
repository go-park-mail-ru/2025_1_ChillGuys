CREATE UNIQUE INDEX IF NOT EXISTS idx_address_unique
    ON bazaar.address (address_string, coordinate);

CREATE UNIQUE INDEX IF NOT EXISTS user_address_unique
    ON bazaar.user_address (user_id, address_id);
