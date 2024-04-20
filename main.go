package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/matmazurk/acc2/db"
	lhttp "github.com/matmazurk/acc2/http"
)

const listenAddr = ":80"

func main() {
	fmt.Println("starting...")
	db, err := db.New("exps.db")
	if err != nil {
		panic(err)
	}
	mux := lhttp.NewMux(db)
	log.Println("listening on", listenAddr)
	err = http.ListenAndServe(listenAddr, mux)
	if err != nil {
		panic(err)
	}
}
