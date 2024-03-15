-- ALTER TABLE products
-- 	DROP CONSTRAINT IF EXISTS product_uid_unique UNIQUE (uid);
-- ALTER TABLE bank_accounts
-- 	DROP CONSTRAINT IF EXISTS bank_account_uid_unique UNIQUE (uid);

-- ALTER TABLE user_transactions
-- 	ADD CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE;
-- ALTER TABLE user_transactions
-- 	ADD CONSTRAINT fk_bank_account_id FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(id) ON DELETE CASCADE;

-- ALTER TABLE user_transactions ALTER COLUMN product_id INT;
-- ALTER TABLE user_transactions ALTER COLUMN bank_account_id INT;
