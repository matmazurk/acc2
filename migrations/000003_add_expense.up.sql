CREATE TABLE IF NOT EXISTS expense (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	category_id INTEGER, 
	payer_id INTEGER,
	amount TEXT NOT NULL,
	currency TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (category_id) REFERENCES category(id),
	FOREIGN KEY (payer_id) REFERENCES payer(id)
);


