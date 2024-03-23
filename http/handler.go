package http

import (
	"embed"
	"html/template"
	"net/http"
	"time"
)

//go:embed templates/*.html
var content embed.FS

type handler struct {
	templates *template.Template
}

func NewMux() *http.ServeMux {
	templates, err := template.ParseFS(content, "templates/*.html")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	h := handler{
		templates: templates,
	}
	h.routes(mux)

	return mux
}

func (h handler) routes(m *http.ServeMux) {
	m.Handle("GET /src/", h.mountSrc())
	m.HandleFunc("GET /", h.getIndex())
	m.HandleFunc("GET /categories", h.getCategories())
	m.HandleFunc("GET /add", h.getAddExpense())
	m.Handle("POST /expenses/add", http.RedirectHandler("/", http.StatusFound))
	m.Handle("POST /categories/add", http.RedirectHandler("/", http.StatusFound))
}

func (h handler) mountSrc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Expires", time.Unix(0, 0).Format(time.RFC1123))

		http.StripPrefix("/src/", http.FileServer(http.Dir("./http/src/"))).ServeHTTP(w, r)
	})
}

func (h handler) getIndex() http.HandlerFunc {
	type Expense struct {
		Description string
		Person      string
		Amount      string
		Currency    string
		Time        string
	}
	type data struct {
		Expenses []Expense
	}
	d := data{
		Expenses: []Expense{
			{
				Description: "zakupy biedra",
				Person:      "mat",
				Amount:      "11.23",
				Currency:    "zł",
				Time:        "23-03-2024 13:33",
			},
			{
				Description: "wazne wydatki",
				Person:      "mat",
				Amount:      "322.43",
				Currency:    "€",
				Time:        "24-03-2024 14:33",
			},
			{
				Description: "dupsko",
				Person:      "mat",
				Amount:      "32.43",
				Currency:    "zł",
				Time:        "22-03-2024 14:33",
			},
			{
				Description: "dlugi opis zakupuw dupa oko sklep",
				Person:      "mat",
				Amount:      "32.43",
				Currency:    "zł",
				Time:        "14:33 22-03-2024",
			},
		},
	}
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			h.templates.ExecuteTemplate(w, "index.html", d)
		})
}

func (h handler) getCategories() http.HandlerFunc {
	type Category string
	type ddata struct {
		Categories []Category
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := ddata{
			Categories: []Category{
				"jedzenie", "zakupy", "zachcianki",
			},
		}

		h.templates.ExecuteTemplate(w, "categories.html", data)
	})
}

func (h handler) getAddExpense() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Data struct {
			Users      []string
			Categories []string
		}
		data := Data{
			Users:      []string{"mat", "paulka"},
			Categories: []string{"zakupy", "zachcianki"},
		}
		h.templates.ExecuteTemplate(w, "add.html", data)
	})
}
