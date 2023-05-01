package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptoRand "crypto/rand"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/TruFaaS/TruFaaS/constants"
	"github.com/TruFaaS/TruFaaS/fission"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	// Access the variables from the .env file
	noOfRuns := 10
	noOfFunctions := []int{10, 25, 50, 100, 250, 500, 750, 1000, 1500, 2000, 2500}
	apiURL := "http://localhost:8080"

	fnInfo := readInfo()

	// Generate a private key
	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), cryptoRand.Reader)

	// Get the public key from the private key
	pubKey := privKey.PublicKey

	pubKeyBytes := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)
	pubKeyHex := hex.EncodeToString(pubKeyBytes)

	// Loop through the values in the slice
	for _, f := range noOfFunctions {
		for i := 0; i < noOfRuns; i++ {

			fmt.Println("<<<<<<<<<<<<<<<< Run ", i+1, " >>>>>>>>>>>>>>>>>>>")

			createInitFunctions(f, &fnInfo, apiURL)

			//create 'f'th function
			creationTime := createTestFunction(&fnInfo, apiURL)

			//run 'f'th function
			runTime := runTestFunction(&fnInfo, apiURL, pubKeyHex)
			cleanUp(apiURL)

			writeToCSV(f, i, creationTime, runTime)

		}
	}

}

func readInfo() fission.Function {
	// Read the JSON file
	file, err := os.ReadFile("info.json")
	if err != nil {
		fmt.Println("Error reading file:", err)
		panic(err)
	}

	// Convert the JSON data to a struct
	var fnInfo fission.Function
	err = json.Unmarshal(file, &fnInfo)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		panic(err)
	}

	return fnInfo
}

func createInitFunctions(f int, fnInfo *fission.Function, apiURL string) {

	for j := 0; j < f-1; j++ {
		fnInfo.FunctionInformation.Name = generateRandomString(5)
		// Send an HTTP POST request with the JSON data in the request body
		jsonData, err := json.Marshal(fnInfo)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			panic(err)
		}
		resp, err := http.Post(apiURL+"/fn/create", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error sending request:", err)
			panic(err)
		}
		defer resp.Body.Close()
	}
	fmt.Println(f-1, "dummy functions created")
}

func createTestFunction(fnInfo *fission.Function, apiURL string) int64 {
	fnInfo.FunctionInformation.Name = "test"
	// Send an HTTP POST request with the JSON data in the request body
	jsonData, err := json.Marshal(fnInfo)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		panic(err)
	}
	start := time.Now() // Record the start time
	resp, err := http.Post(apiURL+"/fn/create", "application/json", bytes.NewBuffer(jsonData))
	elapsed := time.Since(start)
	if err != nil {
		fmt.Println("Error sending request:", err)
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Test function created")
	return elapsed.Microseconds()
}

func runTestFunction(fnInfo *fission.Function, apiURL string, pubKeyHex string) int64 {
	fnInfo.FunctionInformation.Name = "test"
	// Send an HTTP POST request with the JSON data in the request body
	jsonData, err := json.Marshal(fnInfo)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		panic(err)
	}
	// Create a new request with the appropriate method, URL, and body
	req, err := http.NewRequest("POST", apiURL+"/fn/verify", bytes.NewBuffer(jsonData))
	if err != nil {
		// handle error
	}

	// Set the Content-Type header to "application/json"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(constants.InvokerPublicKeyHeader, pubKeyHex)

	start := time.Now() // Record the start time
	// Send the request and get the response
	resp, err := http.DefaultClient.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Error sending request:", err)
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("Test function created")
	return elapsed.Microseconds()
}

func cleanUp(apiURL string) {

	// Clean the API
	req, err := http.NewRequest("GET", apiURL+"/reset-tree", nil)
	if err != nil {
		panic(err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println("API cleaned")
}

func generateRandomString(length int) string {
	// Define the character set to use for the random string
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRTUVWXYZ1234567890"

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Create a byte slice of the specified length
	randomBytes := make([]byte, length)

	// Fill the byte slice with random characters from the charset
	for i := 0; i < length; i++ {
		randomBytes[i] = charset[rand.Intn(len(charset))]
	}

	// Return the random string
	return string(randomBytes)
}

func writeToCSV(noOfFunctions int, runNumber int, createTime int64, runTime int64) {
	// Check if the file "results.csv" exists.
	_, err := os.Stat("results.csv")
	var file *os.File
	var writer *csv.Writer
	if os.IsNotExist(err) {
		// The file does not exist, so create a new file.
		file, err = os.Create("results.csv")
		if err != nil {
			panic(err)
		}

		// Create a CSV writer.
		writer = csv.NewWriter(file)

		// Write the headers to the file.
		headers := []string{"Function Count", "Run Number", "Creation Time", "Run Time"}
		err = writer.Write(headers)
		if err != nil {
			panic(err)
		}
	} else {
		// The file already exists, so open the file in append mode.
		file, err = os.OpenFile("results.csv", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			panic(err)
		}

		// Create a CSV writer.
		writer = csv.NewWriter(file)
	}

	defer file.Close()

	// Write the data to the file.
	data := []string{strconv.Itoa(noOfFunctions), strconv.Itoa(runNumber), strconv.FormatInt(createTime, 10), strconv.FormatInt(runTime, 10)}
	err = writer.Write(data)
	if err != nil {
		panic(err)
	}
	fmt.Println("results written to csv file")
	// Flush the data to the file.
	writer.Flush()
}
