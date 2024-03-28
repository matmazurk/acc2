package http

import (
	"embed"
	"html/template"
	"net/http"
	"time"

	"github.com/matmazurk/acc2/expense"
)

//go:embed templates/*.html
var content embed.FS

type handler struct {
	pers      inter
	templates *template.Template
}

type inter interface {
	Insert(e expense.Expense) error
	SelectExpenses() ([]expense.Expense, error)
	CreatePayer(name string) error
	CreateCategory(name string) error
	ListPayers() ([]string, error)
	ListCategories() ([]string, error)
}

func NewMux(i inter) *http.ServeMux {
	templates, err := template.ParseFS(content, "templates/*.html")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	i.CreatePayer("mat")
	i.CreatePayer("paulka")
	h := handler{
		pers:      i,
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
	m.Handle("POST /expenses/add", h.addExpense())
	m.Handle("POST /categories/add", h.addCategory())
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
		Category    string
		Currency    string
		Time        string
	}
	type data struct {
		Expenses []Expense
	}
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			exps, err := h.pers.SelectExpenses()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			d := data{
				Expenses: make([]Expense, len(exps)),
			}
			for i, e := range exps {
				d.Expenses[i] = Expense{
					Description: e.Description(),
					Person:      e.Payer(),
					Amount:      e.Amount(),
					Category:    e.Category(),
					Currency:    e.Currency(),
					Time:        e.Time().Format("02 Jan 06 15:04"),
				}
			}
			h.templates.ExecuteTemplate(w, "index.html", d)
		})
}

func (h handler) getCategories() http.HandlerFunc {
	type ddata struct {
		Categories []string
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		categories, err := h.pers.ListCategories()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		data := ddata{Categories: categories}
		h.templates.ExecuteTemplate(w, "categories.html", data)
	})
}

func (h handler) getAddExpense() http.HandlerFunc {
	type Data struct {
		Users      []string
		Categories []string
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payers, err := h.pers.ListPayers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		categories, err := h.pers.ListCategories()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		data := Data{
			Users:      payers,
			Categories: categories,
		}
		h.templates.ExecuteTemplate(w, "add.html", data)
	})
}

func (h handler) addExpense() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		description := r.FormValue("description")
		amount := r.FormValue("amount")
		currency := r.FormValue("currency")
		payer := r.FormValue("author")
		category := r.FormValue("category")

		exp, err := expense.NewExpense(description, payer, category, amount, currency, time.Now())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		err = h.pers.Insert(exp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func (h handler) addCategory() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		category := r.FormValue("category")

		err = h.pers.CreateCategory(category)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}
