package main

import (
	"log"
	"net/http"

	"github.com/matmazurk/acc2/db"
	lhttp "github.com/matmazurk/acc2/http"
)

const listenAddr = ":80"

func main() {
	db, err := db.New("exps.db")
	if err != nil {
		panic(err)
	}
	mux := lhttp.NewMux(db)
	log.Println("listening on ", listenAddr)
	http.ListenAndServe(listenAddr, mux)
}
