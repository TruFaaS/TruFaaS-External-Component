package main

import (
	"fmt"
	"net/http"
)

func TestFission(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: Fission")
}
