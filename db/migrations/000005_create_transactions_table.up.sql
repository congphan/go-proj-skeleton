BEGIN;

CREATE TABLE IF NOT EXISTS transactions(
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	account_id INTEGER NOT NULL,
	amount FLOAT(2),
	transaction_type VARCHAR (300),
	created_at VARCHAR (300),
	FOREIGN KEY (user_id) REFERENCES users (id),
	FOREIGN KEY (account_id) REFERENCES accounts (id)
);

COMMIT;