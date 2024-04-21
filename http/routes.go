package http

import (
	"errors"
	"mime"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/matmazurk/acc2/model"
)

func (h handler) routes(m *http.ServeMux) {
	m.Handle("GET /src/", h.mountSrc())
	m.HandleFunc("GET /", h.getIndex())
	m.HandleFunc("GET /categories", h.getCategories())
	m.HandleFunc("GET /add", h.getAddExpense())
	m.Handle("POST /expenses/add", logh(h.addExpense()))
	m.Handle("POST /categories/add", logh(h.addCategory()))
}

func (h handler) mountSrc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Expires", time.Unix(0, 0).Format(time.RFC1123))

		http.StripPrefix("/src/", http.FileServer(http.Dir("./http/src/"))).ServeHTTP(w, r)
	})
}

func (h handler) getIndex() http.HandlerFunc {
	type expense struct {
		Description string
		Person      string
		Amount      string
		Category    string
		Currency    string
		Time        string
	}
	type data struct {
		Expenses []expense
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
				Expenses: make([]expense, len(exps)),
			}
			for i, e := range exps {
				d.Expenses[i] = expense{
					Description: e.Description(),
					Person:      e.Payer(),
					Amount:      e.Amount(),
					Category:    e.Category(),
					Currency:    e.Currency(),
					Time:        e.CreatedAt().Format("02 Jan 06 15:04"),
				}
			}
			h.templates.ExecuteTemplate(w, "index.html", d)
		})
}

func (h handler) getCategories() http.HandlerFunc {
	type data struct {
		Categories []string
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		categories, err := h.pers.ListCategories()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		data := data{Categories: categories}
		h.templates.ExecuteTemplate(w, "categories.html", data)
	})
}

func (h handler) getAddExpense() http.HandlerFunc {
	type data struct {
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

		data := data{
			Users:      payers,
			Categories: categories,
		}
		h.templates.ExecuteTemplate(w, "add.html", data)
	})
}

func (h handler) addExpense() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
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

		exp, err := model.ExpenseBuilder{
			Description: description,
			Payer:       payer,
			Category:    category,
			Amount:      amount,
			Currency:    currency,
			CreatedAt:   time.Now(),
		}.Build()
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

func (h handler) savePhoto(r *http.Request, e model.Expense) error {
	file, header, err := r.FormFile("photo")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil
		}

		return err
	}
	defer file.Close()
	contentType := header.Header.Get("Content-Type")

	// Extract file extension from content type
	exts, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return err
	}

	var ext string
	if slices.Contains(exts, ".jpeg") {
		ext = "jpeg"
	} else if len(exts) > 0 {
		// Get the first extension
		ext = exts[0]
	} else {
		// Fallback to extracting extension from filename
		ext = extractExtension(header.Filename)
	}

	return h.store.SaveExpensePhoto(e, ext, file)
}

func extractExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
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