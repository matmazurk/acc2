package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/matmazurk/acc2/model"
	"github.com/pkg/errors"
)

type db struct {
	db *gorm.DB
}

func New(path string) (db, error) {
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return db{}, errors.Errorf("could not connect to database under '%s'", path)
	}

	err = database.AutoMigrate(&expense{})
	if err != nil {
		return db{}, errors.Wrap(err, "could not auto migrate expense table")
	}

	return db{db: database}, nil
}

func (d db) Insert(e model.Expense) error {
	var p payer
	res := d.db.First(&p, "name = ?", e.Payer())
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.Errorf("payer '%s' not found", e.Payer())
		}
		return res.Error
	}

	var c category
	res = d.db.First(&c, "name = ?", e.Category())
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return errors.Errorf("category '%s' not found", e.Payer())
		}
		return res.Error
	}

	toInsert := expense{
		ID:          e.ID(),
		Description: e.Description(),
		PayerID:     p.ID,
		CategoryID:  c.ID,
		Amount:      e.Amount(),
		Currency:    e.Currency(),
		CreatedAt:   e.Time(),
	}
	res = d.db.Create(&toInsert)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (d db) SelectExpenses() ([]model.Expense, error) {
	var exps []expense
	res := d.db.Preload("Payer").Preload("Category").Model(&expense{}).Order("created_at DESC").Find(&exps)
	if res.Error != nil {
		return nil, res.Error
	}

	es := make([]model.Expense, len(exps))
	var err error
	for i, e := range exps {
		es[i], err = e.toExpense()
		if err != nil {
			return nil, err
		}
	}

	return es, nil
}

func (d db) CreatePayer(name string) error {
	res := d.db.Create(&payer{Name: name})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (d db) CreateCategory(name string) error {
	res := d.db.Create(&category{Name: name})
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (d db) ListPayers() ([]string, error) {
	var payers []payer
	res := d.db.Find(&payers)
	if res.Error != nil {
		return nil, res.Error
	}

	ret := make([]string, len(payers))
	for i, p := range payers {
		ret[i] = p.Name
	}

	return ret, nil
}

func (d db) ListCategories() ([]string, error) {
	var categories []category
	res := d.db.Find(&categories)
	if res.Error != nil {
		return nil, res.Error
	}

	ret := make([]string, len(categories))
	for i, c := range categories {
		ret[i] = c.Name
	}

	return ret, nil
}
