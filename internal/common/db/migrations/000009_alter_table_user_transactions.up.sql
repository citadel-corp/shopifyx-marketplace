-- ALTER TABLE user_transactions
-- 	DROP CONSTRAINT IF EXISTS fk_product_id;
-- ALTER TABLE user_transactions
-- 	DROP CONSTRAINT IF EXISTS fk_bank_account_id;

-- ALTER TABLE products
-- 	DROP CONSTRAINT IF EXISTS product_uid_unique;
-- ALTER TABLE bank_accounts
-- 	DROP CONSTRAINT IF EXISTS bank_account_uid_unique;

-- ALTER TABLE products
-- 	ADD CONSTRAINT product_uid_unique UNIQUE (uid);
-- ALTER TABLE bank_accounts
-- 	ADD CONSTRAINT bank_account_uid_unique UNIQUE (uid);

-- ALTER TABLE user_transactions
-- 	ADD CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES products(uid) ON DELETE CASCADE;
-- ALTER TABLE user_transactions
-- 	ADD CONSTRAINT fk_bank_account_id FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(uid) ON DELETE CASCADE;

-- ALTER TABLE user_transactions ALTER COLUMN product_id SET DATA TYPE UUID USING (gen_random_uuid());
-- ALTER TABLE user_transactions ALTER COLUMN bank_account_id SET DATA TYPE UUID USING (gen_random_uuid());
