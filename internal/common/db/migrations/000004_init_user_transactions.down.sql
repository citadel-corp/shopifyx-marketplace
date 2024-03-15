DROP INDEX IF EXISTS user_transactions_user_id;
DROP INDEX IF EXISTS user_transactions_user_id_product_id;

ALTER TABLE products
	DROP CONSTRAINT product_uid_unique;
ALTER TABLE bank_accounts
	DROP CONSTRAINT bank_account_uid_unique;

ALTER TABLE user_transactions
	DROP CONSTRAINT fk_user_id;
ALTER TABLE user_transactions
	DROP CONSTRAINT fk_product_id;
ALTER TABLE user_transactions
	DROP CONSTRAINT fk_bank_account_id;

DROP TABLE IF EXISTS user_transactions;
