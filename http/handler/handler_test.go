package handler_test

import (
	"io"
	"testing"

	ihttp "github.com/matmazurk/acc2/http"
	"github.com/matmazurk/acc2/model"
	"github.com/rs/zerolog"
)

func TestAddExpense(t *testing.T) {
	pf := newPersistenceFake()
	is := newImagestoreFake()
	m := ihttp.NewMux(pf, is, zerolog.New(zerolog.Nop()))
	_ = m
}

type persistenceFake struct {
	expenses   []model.Expense
	payers     []string
	categories []string
}

func newPersistenceFake() *persistenceFake {
	return &persistenceFake{
		expenses:   []model.Expense{},
		payers:     []string{},
		categories: []string{},
	}
}

func (pf *persistenceFake) Insert(e model.Expense) error {
	pf.expenses = append(pf.expenses, e)
	return nil
}

func (pf *persistenceFake) SelectExpenses() ([]model.Expense, error) {
	return pf.expenses, nil
}

func (pf *persistenceFake) CreatePayer(name string) error {
	pf.payers = append(pf.payers, name)
	return nil
}

func (pf *persistenceFake) CreateCategory(name string) error {
	pf.categories = append(pf.categories, name)
	return nil
}

func (pf *persistenceFake) ListPayers() ([]string, error) {
	return pf.payers, nil
}

func (pf *persistenceFake) ListCategories() ([]string, error) {
	return pf.categories, nil
}

type imagestoreFake struct {
	photos map[string][]byte
}

func newImagestoreFake() *imagestoreFake {
	return &imagestoreFake{
		photos: map[string][]byte{},
	}
}

func (isf *imagestoreFake) SaveExpensePhoto(e model.Expense, fileExtension string, r io.ReadCloser) error {
	contents, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	isf.photos[e.ID()+fileExtension] = contents
	return nil
}

func (isf *imagestoreFake) getPhoto(e model.Expense, fileExtension string) []byte {
	return isf.photos[e.ID()+fileExtension]
}
