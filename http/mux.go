package http

import (
	"net/http"

	"github.com/matmazurk/acc2/http/handler"
)

func NewMux(i handler.Persistence, s handler.Imagestore) *http.ServeMux {
	mux := http.NewServeMux()
	i.CreatePayer("mat")
	i.CreatePayer("paulka")
	h, err := handler.NewHandler(i, s)
	if err != nil {
		panic(err)
	}
	h.Routes(mux)

	return mux
}
