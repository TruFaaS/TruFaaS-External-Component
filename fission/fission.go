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

	fmt.Printf("%#v", mt.Nodes)

	err = utils.StoreMerkleTree(mt)
	if err != nil {
		println(err)
		return
	}

	//send a json response back
	err = utils.SendSuccessResponse(respWriter, http.StatusCreated, nil)
	if err != nil {
		println(err)
		return
	}

}
