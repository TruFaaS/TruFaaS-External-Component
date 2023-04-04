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

//var previousPCRValue []byte

func GetInstanceAtCreate() *simulator.Simulator {
	if sim == nil {
		sim, _ = simulator.Get()
	}
	return sim
}
func SaveToTPM(sim *simulator.Simulator, hashedContent []byte, pcrIndex int) error {

	pcrHandle := tpmutil.Handle(uint32(pcrIndex))

	// uncomment for debugging
	//initialPcrValue, err := tpm2.ReadPCR(sim, pcrIndex, tpm2.AlgSHA256)
	//if err != nil {
	//	log.Fatalf("failed to read PCR: %v", err)
	//}
	err := tpm2.PCRReset(sim, pcrHandle)
	//previousPCR, err := tpm2.ReadPCR(sim, pcrIndex, tpm2.AlgSHA256)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//previousPCRValue = previousPCR

	// TPM PCR extensions follow the calculation:
	// pcr_new = H(pcr_old | H(data))
	// The variable hashedContent already contains the H(data) value
	err = tpm2.PCRExtend(sim, pcrHandle, tpm2.AlgSHA256, hashedContent, "")
	if err != nil {
		fmt.Println(err)
		return err
	}

	pcrValue, err := tpm2.ReadPCR(sim, pcrIndex, tpm2.AlgSHA256)
	if err != nil {
		log.Fatalf("failed to read PCR: %v", err)
	}
	fmt.Println(pcrValue)
	return nil

	//h := sha256.New()
	//h.Write(initialPcrValue)
	//h.Write(hashedContent)
	//hashedValue := h.Sum(nil)
	//
	//// Compare the hash with the value read from the PCR.
	//if bytes.Equal(hashedValue, pcrValue) {
	//	fmt.Println("PCR value matches the extended value.")
	//} else {
	//	fmt.Println("PCR value does not match the extended value.")
	//	fmt.Println(hashedValue)
	//	fmt.Println(pcrValue)
	//}

}

func VerifyMerkleRoot(sim *simulator.Simulator, merkleRoot []byte, pcrIndex int) bool {
	// Read the merkle root stored in the TPM
	pcrValue, err := tpm2.ReadPCR(sim, pcrIndex, tpm2.AlgSHA256)
	if err != nil {
		log.Fatalf("failed to read PCR: %v", err)
		return false
	}
	// This method involves manually recreating the PCRExtend operation
	// TPM PCR extensions follow the calculation:
	// pcr_new = H(pcr_old | H(data))
	// The variable hashedContent already contains the H(data) value
	// Creating a byte array of 32 0's
	// TODO: replace if the PCR will not be reset
	zeroByteArray := bytes.Repeat([]byte{0}, 32)

	// Get the SHA256 algorithm
	h := merkleTree.NewHashFunc()
	// Write the 0 byte array
	h.Write(zeroByteArray)
	// concatenate the merkle root given from TruFaaS
	h.Write(merkleRoot)
	// calculate the hashed value
	hashedValue := h.Sum(nil)

	if bytes.Equal(hashedValue, pcrValue) {
		fmt.Println("PCR value matches the extended value.")
		return true
	} else {
		fmt.Println("PCR value does not match the extended value.")
		fmt.Println(hashedValue)
		fmt.Println(pcrValue)
		return false
	}

}
