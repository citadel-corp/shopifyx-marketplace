CREATE TABLE IF NOT EXISTS user_transactions (
	id SERIAL PRIMARY KEY,
  user_id INT,
  product_id UUID,
  bank_account_id UUID,
  image_url TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE user_transactions DROP CONSTRAINT IF EXISTS fk_user_id;
ALTER TABLE user_transactions DROP CONSTRAINT IF EXISTS fk_product_id;
ALTER TABLE user_transactions DROP CONSTRAINT IF EXISTS fk_bank_account_id;

ALTER TABLE products
	DROP CONSTRAINT IF EXISTS product_uid_unique;
ALTER TABLE bank_accounts
	DROP CONSTRAINT IF EXISTS bank_account_uid_unique;

ALTER TABLE products
	ADD CONSTRAINT product_uid_unique UNIQUE (uid);
ALTER TABLE bank_accounts
	ADD CONSTRAINT bank_account_uid_unique UNIQUE (uid);

ALTER TABLE user_transactions
	ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE user_transactions
	ADD CONSTRAINT fk_product_id FOREIGN KEY (product_id) REFERENCES products(uid) ON DELETE CASCADE;
ALTER TABLE user_transactions
	ADD CONSTRAINT fk_bank_account_id FOREIGN KEY (bank_account_id) REFERENCES bank_accounts(uid) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS user_transactions_user_id
	ON user_transactions (user_id);
CREATE INDEX IF NOT EXISTS user_transactions_user_id_product_id
	ON user_transactions (user_id, product_id);
    