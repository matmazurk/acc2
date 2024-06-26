package handler

import (
	"embed"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/matmazurk/acc2/model"
)

//go:embed src/*
var src embed.FS

func (h handler) Routes(m *http.ServeMux) {
	m.Handle("GET /src/", h.MountSrc())
	m.HandleFunc("GET /", h.GetIndex())

	m.HandleFunc("GET /categories/add", h.GetCategories())
	m.Handle("POST /categories", logh(h.AddCategory(), h.logger))

	m.HandleFunc("GET /expenses/add", h.GetAddExpense())
	m.Handle("POST /expenses", logh(h.AddExpense(), h.logger))
	m.Handle("POST /expenses/{id}/delete", logh(h.DeleteExpense(), h.logger))
	m.Handle("GET /expenses/{id}/photo", h.GetPhoto())
}

func (h handler) MountSrc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Expires", time.Unix(0, 0).Format(time.RFC1123))

		http.FileServer(http.FS(src)).ServeHTTP(w, r)
	})
}

func (h handler) GetIndex() http.HandlerFunc {
	type expense struct {
		ID          string
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
					ID:          e.ID(),
					Description: e.Description(),
					Person:      e.Payer(),
					Amount:      e.Amount(),
					Category:    e.Category(),
					Currency:    e.Currency(),
					Time:        e.CreatedAt().In(h.location).Format("02 Jan 06 15:04"),
				}
			}
			h.templates.ExecuteTemplate(w, "index.html", d)
		})
}

func (h handler) GetCategories() http.HandlerFunc {
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

func (h handler) GetAddExpense() http.HandlerFunc {
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

func (h handler) AddExpense() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			if errors.Is(err, http.ErrNotMultipart) {
				h.logger.Warn().Err(err).Msg("received request with invalid content type")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			h.logger.Error().Err(err).Msg("could not parse multipart form")
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
			CreatedAt:   time.Now().In(h.location),
		}.Build()
		if err != nil {
			h.logger.Warn().Err(err).Msg("invalid request for adding new expense")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = h.savePhoto(r, exp)
		if err != nil {
			h.logger.Error().Err(err).Msg("could not save photo")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		err = h.pers.Insert(exp)
		if err != nil {
			h.logger.Error().Err(err).Msg("could not insert new expense")
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
		ext = ".jpeg"
	} else if len(exts) > 0 {
		ext = exts[0]
	} else {
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

func (h handler) AddCategory() http.HandlerFunc {
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

func (h handler) DeleteExpense() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")

		exps, err := h.pers.SelectExpenses()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		idx := slices.IndexFunc(exps, func(e model.Expense) bool { return e.ID() == idString })
		if idx == -1 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("expense '" + idString + "' not found"))
			return
		}

		err = h.pers.RemoveExpense(exps[idx])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	})
}

func (h handler) GetPhoto() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")

		exps, err := h.pers.SelectExpenses()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		idx := slices.IndexFunc(exps, func(e model.Expense) bool { return e.ID() == idString })
		if idx == -1 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("expense '" + idString + "' not found"))
			return
		}

		photo, err := h.store.LoadExpensePhoto(exps[idx])
		if err != nil {
			if os.IsNotExist(err) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("expense '" + idString + "' has no photo"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer photo.Close()

		io.Copy(w, photo)
	})
}
