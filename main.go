package main

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

const listenAddr = ":80"

type Category string

type ddata struct {
	Categories []Category
}

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

func main() {
	templ, err := template.ParseFiles("http/templates/index.html")
	if err != nil {
		panic(err)
	}
	add, err := template.ParseFiles("http/templates/add.html")
	if err != nil {
		panic(err)
	}
	categories, err := template.ParseFiles("http/templates/categories.html")
	if err != nil {
		panic(err)
	}

	data := data{
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

	http.Handle("GET /src/", http.StripPrefix("/src/", NoCache(http.FileServer(http.Dir("./src")))))

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		templ.Execute(w, data)
	})
	http.HandleFunc("GET /add", func(w http.ResponseWriter, r *http.Request) {
		type Data struct {
			Users      []string
			Categories []string
		}
		data := Data{
			Users:      []string{"mat", "paulka"},
			Categories: []string{"zakupy", "zachcianki"},
		}
		add.Execute(w, data)
	})
	http.HandleFunc("GET /categories", func(w http.ResponseWriter, r *http.Request) {
		data := ddata{
			Categories: []Category{
				"jedzenie", "zakupy", "zachcianki",
			},
		}

		categories.Execute(w, data)
	})
	http.Handle("POST /expenses/add", http.RedirectHandler("/", http.StatusFound))
	http.Handle("POST /categories/add", http.RedirectHandler("/", http.StatusFound))

	log.Println("listening on ", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}

var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func NoCache(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
