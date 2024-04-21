package handler

import (
	"embed"
	"html/template"
	"io"

	"github.com/matmazurk/acc2/model"
	"github.com/rs/zerolog"
)

//go:embed templates/*.html
var content embed.FS

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
	l zerolog.Logger,
) (handler, error) {
	templates, err := template.ParseFS(content, "templates/*.html")
	if err != nil {
		return handler{}, err
	}
	return handler{
		pers:      p,
		store:     is,
		templates: templates,
		logger:    l,
	}, nil
}
