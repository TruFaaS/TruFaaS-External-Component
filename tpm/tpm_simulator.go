package tpm

import (
	"bytes"
	"fmt"
	merkleTree "github.com/TruFaaS/TruFaaS/merkle_tree"
	"github.com/google/go-tpm-tools/simulator"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
	"log"
)

var sim *simulator.Simulator
var pcrIndex int = 23

//var previousPCRValue []byte

func GetInstance() *simulator.Simulator {
	if sim == nil {
		sim, _ = simulator.Get()
	}
	return sim
}
func SaveToTPM(sim *simulator.Simulator, hashedContent []byte) error {

	pcrHandle := tpmutil.Handle(uint32(pcrIndex))

	// uncomment for debugging
	//initialPcrValue, err := tpm2.ReadPCR(sim, pcrIndex, tpm2.AlgSHA256)
	//if err != nil {
	//	log.Fatalf("failed to read PCR: %v", err)
	//}
	err := tpm2.PCRReset(sim, pcrHandle)
	//previousPCR, err := tpm2.ReadPCR(sim, pcrIndex, tpm2.AlgSHA256)
	if err != nil {
		log.Fatalf("failed to reset PCR: %v", err)
		return err
	}
	//previousPCRValue = previousPCR

	// TPM PCR extensions follow the calculation:
	// pcr_new = H(pcr_old | H(data))
	// The variable hashedContent already contains the H(data) value
	err = tpm2.PCRExtend(sim, pcrHandle, tpm2.AlgSHA256, hashedContent, "")
	if err != nil {
		log.Fatalf("failed to extend PCR: %v", err)
		return err
	}

	return nil

}

func VerifyMerkleRoot(sim *simulator.Simulator, merkleRoot []byte) (bool, []byte) {
	// Read the merkle root stored in the TPM
	pcrValue, err := tpm2.ReadPCR(sim, pcrIndex, tpm2.AlgSHA256)
	if err != nil {
		log.Fatalf("failed to read PCR: %v", err)
		return false, nil
	}
	// This method involves manually recreating the PCRExtend operation
	// TPM PCR extensions follow the calculation:
	// pcr_new = H(pcr_old | H(data))
	// The variable hashedContent already contains the H(data) value
	// Creating a byte array of 32 0's
	zeroByteArray := bytes.Repeat([]byte{0}, 32)

	// Get the SHA256 algorithm
	hashCalculator := merkleTree.NewHashFunc()
	// Write the 0 byte array
	hashCalculator.Write(zeroByteArray)
	// Concatenate the merkle root given from TruFaaS
	hashCalculator.Write(merkleRoot)
	// Calculate the hashed value
	hashedValue := hashCalculator.Sum(nil)

	if bytes.Equal(hashedValue, pcrValue) {
		fmt.Println("PCR value matches the extended value.")
		return true, merkleRoot
	} else {
		fmt.Println("PCR value does not match the extended value.")
		return false, nil
	}

}
