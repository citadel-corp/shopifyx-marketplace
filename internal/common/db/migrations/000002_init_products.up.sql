DROP TYPE IF EXISTS product_condition;
CREATE TYPE product_condition AS ENUM ('new', 'second');

CREATE TABLE IF NOT EXISTS products (
	id SERIAL PRIMARY KEY,
	uid UUID NOT NULL,
	name VARCHAR(60) NOT NULL,
	image_url TEXT NOT NULL,
	stock INT NOT NULL,
	condition product_condition NOT NULL,
	tags text[] NOT NULL,
	is_purchaseable boolean NOT NULL,
	price INT NOT NULL,
	purchase_count INT NOT NULL DEFAULT 0,
	user_id INT NOT NULL,
	created_at TIMESTAMP NOT NULL
);

ALTER TABLE products DROP CONSTRAINT IF EXISTS fk_user_id;

ALTER TABLE products
	ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS products_tags
	ON products USING gin(tags);
CREATE INDEX IF NOT EXISTS products_uid
	ON products USING HASH (uid);
CREATE INDEX IF NOT EXISTS products_condition
	ON products(condition);
CREATE INDEX IF NOT EXISTS products_price
	ON products(price);
CREATE INDEX IF NOT EXISTS products_price_asc
	ON products(price ASC);
CREATE INDEX IF NOT EXISTS products_price_desc
	ON products(price DESC);
CREATE INDEX IF NOT EXISTS products_user_id
	ON products (user_id);
CREATE INDEX IF NOT EXISTS products_created_at
	ON products (created_at);
CREATE INDEX IF NOT EXISTS products_created_at_asc
	ON products(created_at ASC);
CREATE INDEX IF NOT EXISTS products_created_at_desc
	ON products(created_at DESC);
CREATE INDEX IF NOT EXISTS products_stock_show_not_empty_only 
	ON products(stock) where stock > 0;
