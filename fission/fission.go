package fission

import (
	"encoding/json"
	"fmt"
	commonTypes "github.com/TruFaaS/TruFaaS/common_types"
	merkleTree "github.com/TruFaaS/TruFaaS/merkle_tree"
	"github.com/TruFaaS/TruFaaS/tpm"
	"github.com/TruFaaS/TruFaaS/utils"
	"net/http"
)

func CreateFnTrustValue(respWriter http.ResponseWriter, req *http.Request) {

	var function Function
	var mt *merkleTree.MerkleTree
	errResponse := commonTypes.ErrorResponse{}

	// get the json value and convert to struct
	err := json.NewDecoder(req.Body).Decode(&function)
	if err != nil {
		errResponse.StatusCode = http.StatusBadRequest
		errResponse.ErrorMsg = err.Error()
		utils.SendErrorResponse(respWriter, errResponse)
		return
	}
	// retrieves already existing merkle tree
	mt, err = utils.RetrieveMerkleTree()
	if err != nil {
		errResponse.StatusCode = http.StatusInternalServerError
		errResponse.ErrorMsg = "Internal Server error"
		utils.SendErrorResponse(respWriter, errResponse)
		return
	}
	// convert the function to byte[]
	fnByteArr, err := json.Marshal(function)
	if err != nil {
		errResponse.StatusCode = http.StatusInternalServerError
		errResponse.ErrorMsg = "Internal Server error"
		return
	}

	mt = mt.AppendNewContent(fnByteArr)

	err = utils.StoreMerkleTree(mt)
	if err != nil {
		println(err)
		return
	}

	// response body
	responseBody := commonTypes.SuccessResponse{StatusCode: http.StatusCreated, Msg: "Function trust value created successfully", FnName: function.FunctionInformation.Name}
	//send a json response back
	utils.SendSuccessResponse(respWriter, responseBody)

}

func VerifyFnTrustValue(respWriter http.ResponseWriter, req *http.Request) {
	var function Function
	var mt *merkleTree.MerkleTree
	errResponse := commonTypes.ErrorResponse{}

	// get the json value and convert to struct
	err := json.NewDecoder(req.Body).Decode(&function)
	if err != nil {
		errResponse.StatusCode = http.StatusBadRequest
		errResponse.ErrorMsg = err.Error()
		utils.SendErrorResponse(respWriter, errResponse)
		return
	}
	// retrieves already existing merkle tree
	mt, err = utils.RetrieveMerkleTree()
	if err != nil {
		errResponse.StatusCode = http.StatusInternalServerError
		errResponse.ErrorMsg = "Internal Server error"
		utils.SendErrorResponse(respWriter, errResponse)
		return
	}

	rootHash := mt.MerkleRoot()

	fnByteArr, err := json.Marshal(function)
	if err != nil {
		errResponse.StatusCode = http.StatusInternalServerError
		errResponse.ErrorMsg = "Internal Server error"
		return
	}
	tpm.TPMTest()
	isVerified := mt.VerifyContentHash(fnByteArr, rootHash)
	if isVerified {
		successResponse := commonTypes.SuccessResponse{StatusCode: http.StatusOK, Msg: "Function verification is successful", FnName: function.FunctionInformation.Name, TrustVerified: true}
		utils.SendSuccessResponse(respWriter, successResponse)
		fmt.Println("verification successful", function.FunctionInformation.Name)

	} else {
		errResponse.StatusCode = http.StatusBadRequest
		errResponse.ErrorMsg = "Function verification failed"
		errResponse.FnName = function.FunctionInformation.Name
		falseVal := false
		errResponse.TrustVerified = &falseVal
		utils.SendErrorResponse(respWriter, errResponse)
		fmt.Println("verification failed", function.FunctionInformation.Name)

	}

}

//;TODO add logger later
