package tpm

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/google/go-tpm-tools/simulator"
	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpmutil"
)

var sim *simulator.Simulator

func GetInstanceAtCreate() *simulator.Simulator {
	if sim == nil {
		sim, _ = simulator.Get()
	}
	return sim
}
func SaveToTPM(sim *simulator.Simulator, content []byte) {
	content = bytes.Repeat([]byte{0xF}, sha256.Size)
	fmt.Println("================Before writing to TPM===========")
	GetFromTPM(sim)

	pcrHandle := tpmutil.Handle(uint32(16))
	err := tpm2.PCRReset(sim, pcrHandle)
	if err != nil {
		fmt.Println(err)
		return
	}
	//pcrValue := bytes.Repeat([]byte{0xF}, sha256.Size)

	err = tpm2.PCRExtend(sim, pcrHandle, tpm2.AlgSHA256, content, "")
	//err = tpm2.PCREvent(sim, pcrHandle, content)
	if err != nil {
		fmt.Println(err)
		return
	}
	expectedValue := tpmutil.U16Bytes(content[:])

	fmt.Println("===============content=============")
	fmt.Println(expectedValue)
	//fmt.Println(len(content))

	//
	fmt.Println("================After writing to TPM==============")
	GetFromTPM(sim)
	//fmt.Println("================Manual hashing")
	//h := sha256.New()
	//h.Write(content)
	//value := h.Sum(nil)
	//fmt.Println("===========Byte array that was hashed")
	//fmt.Println(value)
	//fmt.Println(hex.EncodeToString(value))
}

func GetFromTPM(sim *simulator.Simulator) {

	pcr, err := tpm2.ReadPCR(sim, 16, tpm2.AlgSHA256)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("===========Byte array from reading PCR 16===============")
	fmt.Println(pcr)

	//rs, err := tpm2.ReadPCRs(sim, tpm2.PCRSelection{
	//	Hash: tpm2.AlgSHA256,
	//	PCRs: []int{23},
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println("===========Byte array from reading PCR Map========================")
	//fmt.Println(rs)

}

//func TPMTest() {
//	sim, err := simulator.Get()
//	if err != nil {
//		log.Fatalf("failed to initialize sim: %v", err)
//	}
//	defer func(sim *simulator.Simulator) {
//		err := sim.Close()
//		if err != nil {
//
//		}
//	}(sim)
//
//	// reads initial PCR value
//	//should give a 32 bit string of 0's
//	pcr, err := tpm2.ReadPCR(sim, 7, tpm2.AlgSHA256)
//	if err != nil {
//		return
//	}
//
//	fmt.Printf("PCR %d value: %x\n", 7, pcr)
//	fmt.Println(len(pcr))
//
//	// merkle tree root
//	// TODO: replace
//	sealedSecret := []byte{180, 62, 62, 60, 193, 42, 73, 38, 4, 48, 163, 67, 240, 116, 35, 151, 125, 172, 172, 200, 140, 175, 141, 215, 94, 181, 12, 165, 44, 146, 178, 188}
//	fmt.Println(sealedSecret)
//
//	//hash := sha256.Sum256(sealedSecret)
//	pcrHandle := tpmutil.Handle(7)
//	err = tpm2.PCRExtend(sim, pcrHandle, tpm2.AlgNull, sealedSecret, "")
//	err = tpm2.PCREvent(sim, pcrHandle, sealedSecret)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	pcr, err = tpm2.ReadPCR(sim, 7, tpm2.AlgSHA256)
//	if err != nil {
//		return
//	}
//
//	fmt.Printf("PCR %d value: %x\n", 7, pcr)
//	fmt.Println(len(pcr))
//
//}
