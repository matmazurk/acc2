package handler

import (
	"embed"
	"html/template"
	"io"
	"time"

	"github.com/matmazurk/acc2/model"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

//go:embed templates/*.html
var content embed.FS

type Persistence interface {
	Insert(e model.Expense) error
	RemoveExpense(e model.Expense) error
	SelectExpenses() ([]model.Expense, error)
	CreatePayer(name string) error
	CreateCategory(name string) error
	ListPayers() ([]string, error)
	ListCategories() ([]string, error)
}

type Imagestore interface {
	SaveExpensePhoto(e model.Expense, fileExtension string, r io.ReadCloser) error
	LoadExpensePhoto(e model.Expense) (io.ReadCloser, error)
}

type handler struct {
	pers      Persistence
	store     Imagestore
	templates *template.Template
	location  *time.Location
	logger    zerolog.Logger
}

func NewHandler(
	p Persistence,
	is Imagestore,
) (handler, error) {
	templates, err := template.ParseFS(content, "templates/*.html")
	if err != nil {
		return handler{}, err
	}
	loc, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		return handler{}, errors.Wrap(err, "could not load Europe/Warsaw location")
	}
	return handler{
		pers:      p,
		store:     is,
		templates: templates,
		location:  loc,
	}, nil
}
