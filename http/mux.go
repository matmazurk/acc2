package http

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/matmazurk/acc2/http/handler"
	"github.com/rs/zerolog"
)

//go:embed templates/*.html
var content embed.FS

func NewMux(i handler.Persistence, s handler.Imagestore, logger zerolog.Logger) *http.ServeMux {
	templates, err := template.ParseFS(content, "templates/*.html")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	i.CreatePayer("mat")
	i.CreatePayer("paulka")
	h := handler.NewHandler(i, s, templates, logger)
	h.Routes(mux)

	return mux
}
