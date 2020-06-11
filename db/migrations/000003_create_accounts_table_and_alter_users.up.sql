BEGIN;

CREATE TABLE IF NOT EXISTS accounts(
	id INTEGER PRIMARY KEY,
	user_id INTEGER NOT NULL,
	name VARCHAR (300),
	bank VARCHAR (300),
	FOREIGN KEY (user_id) REFERENCES users (id)
);

INSERT INTO accounts (id, user_id, name, bank)
VALUES (1, 1, 'Alice', 'VCB');

INSERT INTO accounts (id, user_id, name, bank)
VALUES (2, 1, 'Alice', 'VIB');

COMMIT;