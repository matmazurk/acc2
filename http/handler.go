package http

import (
	"embed"
	"html/template"
	"io"
	"net/http"

	"github.com/matmazurk/acc2/model"
)

//go:embed templates/*.html
var content embed.FS

type handler struct {
	pers      inter
	store     store
	templates *template.Template
}

type inter interface {
	Insert(e model.Expense) error
	SelectExpenses() ([]model.Expense, error)
	CreatePayer(name string) error
	CreateCategory(name string) error
	ListPayers() ([]string, error)
	ListCategories() ([]string, error)
}

type store interface {
	SaveExpensePhoto(e model.Expense, fileExtension string, r io.ReadCloser) error
}

func NewMux(i inter, s store) *http.ServeMux {
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
	}
	h.routes(mux)

	return mux
}
