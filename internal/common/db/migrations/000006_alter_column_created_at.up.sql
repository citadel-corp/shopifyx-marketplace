ALTER TABLE users ALTER COLUMN created_at SET DEFAULT current_timestamp;
ALTER TABLE products ALTER COLUMN created_at SET DEFAULT current_timestamp;
ALTER TABLE bank_accounts ADD COLUMN created_at timestamp DEFAULT current_timestamp;
ALTER TABLE user_transactions ALTER COLUMN created_at SET DEFAULT current_timestamp;
