package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/matmazurk/acc2/model"
	"github.com/pkg/errors"
	_ "modernc.org/sqlite"
)

type Client struct {
	db *sqlx.DB
}

func New(path string) (Client, error) {
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return Client{}, fmt.Errorf("could not open database: %w", err)
	}

	err = migrateUp(db)
	if err != nil {
		return Client{}, fmt.Errorf("could not migrate up: %w", err)
	}

	return Client{db: db}, nil
}

func (d Client) Insert(e model.Expense) error {
	p, err := d.getPayer(e.Payer())
	if err != nil {
		return err
	}

	c, err := d.getCategory(e.Category())
	if err != nil {
		return err
	}

	_, err = d.db.Exec(`
		INSERT INTO expense(id, category_id, payer_id, amount, currency, description, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.ID(), c.ID, p.ID, e.Amount(), e.Currency(), e.Description(), e.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf("could not insert new expense: %w", err)
	}

	return nil
}

func (d Client) SelectExpenses() ([]model.Expense, error) {
	var exps []expense
	err := d.db.Select(&exps, "SELECT * FROM expenses ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("could not select expenses: %w", err)
	}

	es := make([]model.Expense, len(exps))
	for i, e := range exps {
		es[i], err = model.ExpenseBuilder{
			Id:          e.ID,
			Description: e.Description,
			Payer:       e.Payer.Name,
			Category:    e.CategoryID.Name,
			Amount:      e.Amount,
			Currency:    e.Currency,
			CreatedAt:   e.CreatedAt,
		}.Build()
	}

	return es, nil
}

func (d Client) CreatePayer(name string) error {
	_, err := d.db.Exec("INSERT INTO payer(name) VALUES (?)", name)
	return err
}

func (d Client) CreateCategory(name string) error {
	_, err := d.db.Exec("INSERT INTO category(name) VALUES (?)", name)
	return err
}

func (d Client) ListPayers() ([]string, error) {
	var payers []payer
	err := d.db.Select(&payers, "SELECT * FROM payer")
	if err != nil {
		return nil, fmt.Errorf("could not get all payers: %w", err)
	}

	ret := make([]string, len(payers))
	for i, p := range payers {
		ret[i] = p.Name
	}

	return ret, nil
}

func (d Client) ListCategories() ([]string, error) {
	var categories []category
	err := d.db.Select(&categories, "SELECT * FROM category")
	if err != nil {
		return nil, fmt.Errorf("could not get all categories: %w", err)
	}

	ret := make([]string, len(categories))
	for i, c := range categories {
		ret[i] = c.Name
	}

	return ret, nil
}

func (d Client) RemoveExpense(e model.Expense) error {
	_, err := d.db.Exec("DELETE FROM expense WHERE id = ?", e.ID())
	if err != nil {
		return fmt.Errorf("could not remove expense: %w", err)
	}

	return nil
}

func (d Client) getPayer(name string) (payer, error) {
	var p payer
	err := d.db.Get(&p, "SELECT * FROM payer WHERE name = ?", name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, fmt.Errorf("no such payer: '%s'", name)
		}
		return p, fmt.Errorf("could not get payer: %w", err)
	}

	return p, nil
}

func (d Client) getCategory(name string) (category, error) {
	var c category
	err := d.db.Get(&c, "SELECT * FROM category WHERE name = ?", name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c, fmt.Errorf("no such category: '%s'", name)
		}
		return c, fmt.Errorf("could not get category: %w", err)
	}

	return c, nil
}
