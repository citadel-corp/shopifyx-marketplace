ALTER TABLE users ALTER COLUMN created_at SET DEFAULT;
ALTER TABLE products ALTER COLUMN created_at SET DEFAULT;
ALTER TABLE bank_accounts DROP COLUMN created_at;
ALTER TABLE user_transactions ALTER COLUMN created_at SET DEFAULT;
