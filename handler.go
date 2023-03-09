package main

import (
	"fmt"
	"net/http"
)

func TestFn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: ")
}
