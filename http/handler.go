package http

import (
	"embed"
	"html/template"
	"io"
	"net/http"

	"github.com/matmazurk/acc2/model"
	"github.com/rs/zerolog"
)

//go:embed templates/*.html
var content embed.FS

type handler struct {
	pers      persistence
	store     imagestore
	templates *template.Template
	logger    zerolog.Logger
}

type persistence interface {
	Insert(e model.Expense) error
	SelectExpenses() ([]model.Expense, error)
	CreatePayer(name string) error
	CreateCategory(name string) error
	ListPayers() ([]string, error)
	ListCategories() ([]string, error)
}

type imagestore interface {
	SaveExpensePhoto(e model.Expense, fileExtension string, r io.ReadCloser) error
}

func NewMux(i persistence, s imagestore, logger zerolog.Logger) *http.ServeMux {
	templates, err := template.ParseFS(content, "templates/*.html")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	i.CreatePayer("mat")
	i.CreatePayer("paulka")
	h := handler{
		pers:      i,
		store:     s,
		templates: templates,
		logger:    logger,
	}
	h.routes(mux)

	return mux
}
