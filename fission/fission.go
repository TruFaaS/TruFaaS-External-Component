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
	sim := tpm.GetInstanceAtCreate()
	err = tpm.SaveToTPM(sim, mt.MerkleRoot())
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

	sim := tpm.GetInstanceAtCreate()
	verified, merkleRoot := tpm.VerifyMerkleRoot(sim, rootHash)
	if !verified {
		// TODO: update with appropriate status code
		errResponse.StatusCode = http.StatusInternalServerError
		errResponse.ErrorMsg = "TPM Value and Merkle Root don't match"
		return
	}
	fnByteArr, err := json.Marshal(function)
	if err != nil {
		errResponse.StatusCode = http.StatusInternalServerError
		errResponse.ErrorMsg = "Internal Server error"
		return
	}
	isVerified := mt.VerifyContentHash(fnByteArr, merkleRoot)
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

func TestTPMMethod(respWriter http.ResponseWriter, req *http.Request) {
	sim := tpm.GetInstanceAtCreate()
	sealedSecret := []byte{180, 62, 62, 60, 193, 42, 73, 38, 4, 48, 163, 67, 240, 116, 35, 151, 125, 172, 172, 200, 140, 175, 141, 215, 94, 181, 12, 165, 44, 146, 178, 188}
	//sealedSecret := []byte{1, 2, 3}
	err := tpm.SaveToTPM(sim, sealedSecret)
	if err != nil {
		fmt.Println(err)
		return
	}
	tpm.VerifyMerkleRoot(sim, sealedSecret)
}

//;TODO add logger later
