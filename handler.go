package main

import (
	"fmt"
	"net/http"
)

func TestFission(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: Fission")
}

func TestOF(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: OpenFaaS")
}
