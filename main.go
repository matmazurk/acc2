package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/matmazurk/acc2/db"
	lhttp "github.com/matmazurk/acc2/http"
	"github.com/matmazurk/acc2/imagestore"
)

const listenAddr = ":80"

func main() {
	fmt.Println("starting...")
	db, err := db.New("exps.db")
	if err != nil {
		panic(err)
	}
	store, err := imagestore.NewStore(".")
	if err != nil {
		panic(err)
	}
	mux := lhttp.NewMux(db, store)
	log.Println("listening on", listenAddr)
	err = http.ListenAndServe(listenAddr, mux)
	if err != nil {
		panic(err)
	}
}
