package tpm

import (
	"fmt"
	"github.com/google/go-tpm/tpm2"
)

func main() {
	//sim, err := simulator.Get()
	tpm, err := tpm2.OpenTPM("/home/gayangi/Downloads/Simulator.exe")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(tpm)
	//if err != nil {
	//	log.Fatalf("failed to initialize sim: %v", err)
	//}
	//defer func(sim *simulator.Simulator) {
	//	err := sim.Close()
	//	if err != nil {
	//
	//	}
	//}(sim)
	//
	//// reads initial PCR value
	////should give a 32 bit string of 0's
	//pcr, err := tpm2.ReadPCR(sim, 7, tpm2.AlgSHA256)
	//if err != nil {
	//	return
	//}
	//
	//fmt.Printf("PCR %d value: %x\n", 7, pcr)
	//fmt.Println(len(pcr))
	//
	//// merkle tree root
	//// TODO: replace
	//sealedSecret := []byte{180, 62, 62, 60, 193, 42, 73, 38, 4, 48, 163, 67, 240, 116, 35, 151, 125, 172, 172, 200, 140, 175, 141, 215, 94, 181, 12, 165, 44, 146, 178, 188}
	//fmt.Println(sealedSecret)
	//
	////hash := sha256.Sum256(sealedSecret)
	//pcrHandle := tpmutil.Handle(7)
	//err = tpm2.PCRExtend(sim, pcrHandle, tpm2.AlgNull, sealedSecret, "")
	//err = tpm2.PCREvent(sim, pcrHandle, sealedSecret)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//pcr, err = tpm2.ReadPCR(sim, 7, tpm2.AlgSHA256)
	//if err != nil {
	//	return
	//}
	//
	//fmt.Printf("PCR %d value: %x\n", 7, pcr)
	//fmt.Println(len(pcr))

	//sim, err := simulator.Get()
	//var tpm, err = tpm2.OpenTPM("/home/gayangi/Downloads/IBM/src/tpm_server")
	//conn, err := net.Dial("tcp", "localhost:2321")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//if err != nil {
	//	// handle error
	//}
	//defer conn.Close()

	//TODO - Use as necessary

	// Create a new TPM client using the simulator connection
	//sim, err := tpm2.OpenTPM("/home/gayangi/Downloads/IBM/src/tpm_server")
	//sim, err := tpmutil.OpenTPM()
	////tpm :=tpmutil.OpenTPM()
	////sim, err := tpm2.OpenTPM("tcp://127.0.0.1:2321")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(sim)
	//conn, err := net.Dial("tcp", "127.0.0.1:2321")
	//if err != nil {
	//	fmt.Println("Connection error")
	//	fmt.Println(err)
	//	return
	//}
	//register := 16
	//actualTPM, _ := tpm2.OpenTPM()
	//tpmPath := flag.String()
	//pcr, err := tpm2.ReadPCR(conn, register, tpm2.AlgSHA256)
	//if err != nil {
	//	fmt.Println("Reading error")
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(pcr)
	//content := bytes.Repeat([]byte{0xF}, sha256.Size)
	//err = tpm2.PCRExtend(conn, tpmutil.Handle(register), tpm2.AlgSHA256, content, "")
	//if err != nil {
	//	fmt.Println("write error")
	//	fmt.Println(err)
	//	return
	//}
	//pcr, err = tpm2.ReadPCR(conn, register, tpm2.AlgSHA256)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(pcr)

	//selection := tpm2.PCRSelection{
	//	Hash: tpm2.AlgSHA3_256,
	//	PCRs: []int{16, 23},
	//}
	//pcr, _ := tpm2.ReadPCRs(conn, selection)

}
