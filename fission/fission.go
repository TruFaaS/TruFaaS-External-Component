package fission

import (
	"encoding/json"
	"fmt"
	merkleTree "github.com/TruFaaS/TruFaaS/merkle_tree"
	"github.com/TruFaaS/TruFaaS/utils"
	"net/http"
)

func CreateFnTrustValue(respWriter http.ResponseWriter, req *http.Request) {

	var function Function
	var mt *merkleTree.MerkleTree

	// get the json value and convert to struct
	err := json.NewDecoder(req.Body).Decode(&function)
	if err != nil {
		fmt.Println(err)
		return
	}
	// retrieves already existing merkle tree
	mt, err = utils.RetrieveMerkleTree()
	if err != nil {
		fmt.Println(err)
		return
	}
	// convert the function to byte[]
	fnByteArr, err := json.Marshal(function)
	if err != nil {
		fmt.Println(err)
		return
	}

	mt = mt.AppendNewContent(fnByteArr)

	//fmt.Printf("%#v", mt.Nodes)

	err = utils.StoreMerkleTree(mt)
	if err != nil {
		println(err)
		return
	}

	//send a json response back
	err = utils.SendSuccessResponse(respWriter, http.StatusCreated, "")
	if err != nil {
		println(err)
		return
	}

}

func VerifyFnTrustValue(respWriter http.ResponseWriter, req *http.Request) {
	var function Function
	var mt *merkleTree.MerkleTree

	// get the json value and convert to struct
	err := json.NewDecoder(req.Body).Decode(&function)
	if err != nil {
		fmt.Println(err)
		return
	}
	// retrieves already existing merkle tree
	mt, err = utils.RetrieveMerkleTree()
	if err != nil {
		fmt.Println(err)
		return
	}

	rootHash := mt.MerkleRoot()

	fnByteArr, err := json.Marshal(function)
	if err != nil {
		fmt.Println(err)
		return
	}

	isVerified := mt.VerifyContentHash(fnByteArr, rootHash)
	if isVerified {
		err = utils.SendSuccessResponse(respWriter, http.StatusOK, "Trust is verified")
		if err != nil {
			return
		}
	} else {
		err = utils.SendSuccessResponse(respWriter, http.StatusOK, "Failed to verify trust")
		if err != nil {
			return
		}
	}
}

//;TODO check and add error responses back to client where necessary
