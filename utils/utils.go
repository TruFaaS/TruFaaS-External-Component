package utils

import (
	"encoding/gob"
	"fmt"
	"github.com/TruFaaS/TruFaaS/constants"
	merkleTree "github.com/TruFaaS/TruFaaS/merkle_tree"
	"os"
)

func StoreMerkleTree(tree *merkleTree.MerkleTree) error {
	fmt.Println("storing")
	file, err := os.Create(constants.TreeStoreFileName)
	if err != nil {
		fmt.Println("Failed to create file, ERROR:", err)
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err = encoder.Encode(tree); err != nil {
		fmt.Println("Failed to encode tree into binary, ERROR:", err.Error())
		return err
	}

	fmt.Println("Merkle Tree saved to", constants.TreeStoreFileName)
	return nil
}

func RetrieveMerkleTree() (*merkleTree.MerkleTree, error) {
	_, err := os.Stat(constants.TreeStoreFileName)
	if os.IsNotExist(err) {
		fmt.Println("file not found")
		return merkleTree.NewTree(), nil

	} else {
		file, err := os.Open(constants.TreeStoreFileName)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return nil, err
		}
		defer file.Close()

		decoder := gob.NewDecoder(file)
		var mt *merkleTree.MerkleTree
		if err := decoder.Decode(&mt); err != nil {
			fmt.Println("Error decoding binary into tree:", err)
			return nil, err
		}
		return mt, nil
	}
}
