package fission

import (
	"encoding/json"
	"fmt"
	merkleTree "github.com/TruFaaS/TruFaaS/merkle_tree"
	"github.com/TruFaaS/TruFaaS/utils"
	"net/http"
)

func Create(respWriter http.ResponseWriter, req *http.Request) {

	var function Function
	var mt *merkleTree.MerkleTree

	// get the json value and convert to struct
	err := json.NewDecoder(req.Body).Decode(&function)
	if err != nil {
		fmt.Println(err)
		return
	}
	// create a node in the merkel tree with the function
	//steps:
	//read the merkle tree and assign it to mt
	//update the merkle tree and store it back
	mt, err = utils.RetrieveMerkleTree()
	if err != nil {
		return
	}
	// convert the function to byte[]
	fnByteArr, err := json.Marshal(function)
	if err != nil {
		return
	}
	mt = mt.AppendNewContent(fnByteArr)

	fmt.Printf("%#v", mt)

	err = utils.StoreMerkleTree(mt)
	if err != nil {
		println(err)
		return
	}

	//send a json response back
	//marshal, err := json.Marshal(function)
	//if err != nil {
	//	return
	//}
	//respWriter.Header().Set("Content-Type", constants.ContentTypeJSON)
	//respWriter.WriteHeader(http.StatusCreated)
	//_, err = respWriter.Write(marshal)
	//if err != nil {
	//	return
	//}

	//fmt.Printf("%#v", function)

}
