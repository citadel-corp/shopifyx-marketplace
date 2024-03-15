ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_unique;
ALTER TABLE users ADD CONSTRAINT users_username_unique UNIQUE (username);