package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/test", TestFn)
	log.Fatal(http.ListenAndServe(":10000", r))

}
