package main

import (
	"fmt"
	"github.com/TruFaaS/TruFaaS/constants"
	"net/http"
	"os"
)

func ResetTree(respWriter http.ResponseWriter, req *http.Request) {

	// path to the binary file
	filePath := constants.TreeStoreFileName

	// delete the binary file
	err := os.Remove(filePath)
	if err != nil {
		fmt.Printf("Error deleting file: %v\n", err)
		return
	}

	fmt.Println("Merkle tree file deleted successfully.")

}
