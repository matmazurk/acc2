CREATE VIEW IF NOT EXISTS expenses AS
SELECT e.id, e.amount, e.currency, e.description, e.created_at, c.id AS "category.id", c.name AS "category.name", p.id AS "payer.id", p.name AS "payer.name"  FROM expense e
JOIN category c ON c.id = e.category_id
JOIN payer p ON p.id = e.payer_id;
