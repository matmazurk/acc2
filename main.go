package main

import (
	"log"
	"net/http"

	lhttp "github.com/matmazurk/acc2/http"
)

const listenAddr = ":80"

func main() {
	mux := lhttp.NewMux()
	log.Println("listening on ", listenAddr)
	http.ListenAndServe(listenAddr, mux)
}
