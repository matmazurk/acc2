package handler

import (
	"html/template"
	"io"

	"github.com/matmazurk/acc2/model"
	"github.com/rs/zerolog"
)

type Persistence interface {
	Insert(e model.Expense) error
	SelectExpenses() ([]model.Expense, error)
	CreatePayer(name string) error
	CreateCategory(name string) error
	ListPayers() ([]string, error)
	ListCategories() ([]string, error)
}

type Imagestore interface {
	SaveExpensePhoto(e model.Expense, fileExtension string, r io.ReadCloser) error
}

type handler struct {
	pers      Persistence
	store     Imagestore
	templates *template.Template
	logger    zerolog.Logger
}

func NewHandler(
	p Persistence,
	is Imagestore,
	temps *template.Template,
	l zerolog.Logger,
) handler {
	return handler{
		pers:      p,
		store:     is,
		templates: temps,
		logger:    l,
	}
}
