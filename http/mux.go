package http

import (
	"net/http"

	"github.com/matmazurk/acc2/http/handler"
	"github.com/rs/zerolog"
)

func NewMux(i handler.Persistence, s handler.Imagestore, logger zerolog.Logger) *http.ServeMux {
	mux := http.NewServeMux()
	i.CreatePayer("mat")
	i.CreatePayer("paulka")
	h, err := handler.NewHandler(i, s, logger)
	if err != nil {
		panic(err)
	}
	h.Routes(mux)

	return mux
}
