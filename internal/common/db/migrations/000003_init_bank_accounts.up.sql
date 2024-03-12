CREATE TABLE IF NOT EXISTS bank_accounts (
	id SERIAL PRIMARY KEY,
	uid UUID NOT NULL,
	name VARCHAR(15),
	account_name VARCHAR(15),
	account_number VARCHAR(15),
	user_id INT
	-- CONSTRAINT fk_user_id
	-- 	FOREIGN KEY(user_id)
	-- 		REFERENCES users(id)
	-- 			ON DELETE CASCADE
	-- 			ON UPDATE NO ACTION
);

ALTER TABLE bank_accounts
	ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX bank_accounts_user_id
	ON bank_accounts (user_id);
CREATE INDEX bank_accounts_uid
	ON bank_accounts USING HASH (uid);
    