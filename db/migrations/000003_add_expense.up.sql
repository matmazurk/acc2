CREATE TABLE IF NOT EXISTS expense (
	id TEXT PRIMARY KEY,
	category_id INTEGER, 
	payer_id INTEGER,
	amount TEXT NOT NULL,
	currency TEXT NOT NULL,
	description TEXT NOT NULL,
        created_at DATETIME NOT NULL,

	FOREIGN KEY (category_id) REFERENCES category(id),
	FOREIGN KEY (payer_id) REFERENCES payer(id)
);


