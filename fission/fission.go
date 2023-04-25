package fission

import (
	"encoding/json"
	"fmt"
	commonTypes "github.com/TruFaaS/TruFaaS/common_types"
	"github.com/TruFaaS/TruFaaS/constants"
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
		utils.SendErrorResponse(respWriter, errResponse)
		return
	}

	mt = mt.AppendNewContent(fnByteArr)

	err = utils.StoreMerkleTree(mt)
	if err != nil {
		println(err)
		return
	}
	sim := tpm.GetInstance()
	err = tpm.SaveToTPM(sim, mt.GetMerkleRoot())
	if err != nil {
		println(err)
		return
	}

	// response body
	responseBody := commonTypes.SuccessResponse{StatusCode: http.StatusCreated, Msg: "Function trust value created successfully", FnName: function.FunctionInformation.Name}
	//send a json response back
	utils.SendSuccessResponse(respWriter, responseBody)
	// logs
	fmt.Println("function created successfully, function Name: ", function.FunctionInformation.Name)

}

func VerifyFnTrustValue(respWriter http.ResponseWriter, req *http.Request) {
	var function Function
	var mt *merkleTree.MerkleTree
	errResponse := commonTypes.ErrorResponse{}

	// check and set client public key value
	clientPubKeyHeader := req.Header.Get(constants.ClientPublicKeyHeader)

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

	rootHash := mt.GetMerkleRoot()

	sim := tpm.GetInstance()
	merkleTreeVerifiedWithTpm, merkleRoot := tpm.VerifyMerkleRoot(sim, rootHash)
	if !merkleTreeVerifiedWithTpm {
		utils.SendVerificationFailureErrorResponse(respWriter, function.FunctionInformation.Name, clientPubKeyHeader)
		fmt.Println("verification failed", function.FunctionInformation.Name)
		return
	}

	fnByteArr, err := json.Marshal(function)
	if err != nil {
		errResponse.StatusCode = http.StatusInternalServerError
		errResponse.ErrorMsg = "Internal Server error"
		utils.SendErrorResponse(respWriter, errResponse)
		return
	}

	contentHashVerified := mt.VerifyContentHash(fnByteArr, merkleRoot)
	if contentHashVerified {
		utils.SendVerificationSuccessResponse(respWriter, function.FunctionInformation.Name, clientPubKeyHeader)
		fmt.Println("verification successful, function name: ", function.FunctionInformation.Name)

	} else {
		utils.SendVerificationFailureErrorResponse(respWriter, function.FunctionInformation.Name, clientPubKeyHeader)
		fmt.Println("verification failed", function.FunctionInformation.Name)

	}

}

//;TODO add logger later
