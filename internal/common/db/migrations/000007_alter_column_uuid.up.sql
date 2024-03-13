ALTER TABLE products ALTER COLUMN uid SET DEFAULT gen_random_uuid();
ALTER TABLE bank_accounts ALTER COLUMN uid SET DEFAULT gen_random_uuid();
