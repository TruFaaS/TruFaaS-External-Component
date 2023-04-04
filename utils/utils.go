package utils

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	commonTypes "github.com/TruFaaS/TruFaaS/common_types"
	"github.com/TruFaaS/TruFaaS/constants"
	merkleTree "github.com/TruFaaS/TruFaaS/merkle_tree"
	"net/http"
	"os"
)

// StoreMerkleTree : to store the updated merkle tree
func StoreMerkleTree(tree *merkleTree.MerkleTree) error {

	file, err := os.Create(constants.TreeStoreFileName)
	if err != nil {
		fmt.Println("Failed to create file, ERROR:", err)
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err = encoder.Encode(tree); err != nil {
		fmt.Println("Failed to encode tree into binary, ERROR:", err.Error())
		return err
	}

	fmt.Println("Merkle Tree saved to", constants.TreeStoreFileName)
	return nil
}

// RetrieveMerkleTree : to retrieve the existing merkle tree, or return a new tree if it doesn't exist
func RetrieveMerkleTree() (*merkleTree.MerkleTree, error) {
	_, err := os.Stat(constants.TreeStoreFileName)
	if os.IsNotExist(err) {
		fmt.Println("No exiting merkle tree found")
		return merkleTree.NewTree(), nil

	} else {
		file, err := os.Open(constants.TreeStoreFileName)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return nil, err
		}
		defer file.Close()

		decoder := gob.NewDecoder(file)
		var mt *merkleTree.MerkleTree
		if err = decoder.Decode(&mt); err != nil {
			fmt.Println("Error decoding binary into tree:", err)
			return nil, err
		}
		return mt, nil
	}
}

// SendSuccessResponse SendResponse : tos send the success response back to the client
func SendSuccessResponse(respWriter http.ResponseWriter, body commonTypes.SuccessResponse) {

	jsonResponse, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("failed to marshal body, Error:%s", err)
		return
	}
	respWriter.Header().Set("Content-Type", constants.ContentTypeJSON)
	respWriter.WriteHeader(body.StatusCode)
	_, err = respWriter.Write(jsonResponse)
	if err != nil {
		fmt.Printf("failed to marshal body, Error:%s", err)
		return
	}

}

func SendErrorResponse(respWriter http.ResponseWriter, body commonTypes.ErrorResponse) {

	jsonResponse, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("failed to marshal body, Error:%s", err)
		return
	}
	respWriter.Header().Set("Content-Type", constants.ContentTypeJSON)
	respWriter.WriteHeader(body.StatusCode)
	_, err = respWriter.Write(jsonResponse)
	if err != nil {
		fmt.Printf("failed to marshal body, Error:%s", err)
		return
	}

}

func SendVerificationFailureErrorResponse(respWriter http.ResponseWriter, fnName string) {

	errResponse := commonTypes.ErrorResponse{}
	//:TODO change status code
	errResponse.StatusCode = http.StatusNotFound
	errResponse.ErrorMsg = "Function verification failed"
	errResponse.FnName = fnName
	falseVal := false
	errResponse.TrustVerified = &falseVal

	SendErrorResponse(respWriter, errResponse)

}
