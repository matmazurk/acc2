package main

import (
	"log"
	"net/http"
)

const listenAddr = ":80"

func main() {
	log.Println("listening on ", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}
