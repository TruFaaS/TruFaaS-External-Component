package merkle_tree

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"
)

// MerkleTree represents a Merkle tree
type MerkleTree struct {
	Root           int       // Index of the root node in the nodes slice
	MerkleRootHash []byte    // Root hash of the Merkle tree
	Leafs          []int     // Indexes of the leaf nodes in the nodes slice
	Nodes          []*Node   // Nodes in the Merkle tree
	HashStrategy   hash.Hash // Hash function used to compute Merkle tree hashes
}

// Node represents a node in the Merkle tree
type Node struct {
	Parent int    // Index of the parent node in the nodes slice
	Left   int    // Index of the left child node in the nodes slice
	Right  int    // Index of the right child node in the nodes slice
	Leaf   bool   // Indicates whether the node is a leaf node
	Dup    bool   // Indicates whether the node is a duplicate leaf node
	Hash   []byte // Hash value of the node
}

// NewTree creates a new Merkle Tree with default hash strategy of SHA256
func NewTree() *MerkleTree {
	t := &MerkleTree{Nodes: make([]*Node, 0)}
	return t
}

// MerkleRoot returns the root hash of the Merkle tree
func (t *MerkleTree) MerkleRoot() []byte {
	return t.MerkleRootHash
}

// hashByteSlice returns the hash of a byte slice using the hash function specified in the Merkle tree
func (t *MerkleTree) hashByteSlice(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)[:]
}

// AppendNewContent builds a new tree with the new content and return the tree
func (t *MerkleTree) AppendNewContent(content []byte) *MerkleTree {

	leaf := &Node{
		Parent: -1,
		Left:   -1,
		Right:  -1,
		Leaf:   true,
		Dup:    false,
		Hash:   t.hashByteSlice(content),
	}
	t.Leafs, t.Nodes = updateLeafsSlice(t.Leafs, t.Nodes, leaf)

	// Recursively build intermediate nodes until the root node is reached
	root, err := buildIntermediate(t.Leafs, t)

	// Create the new tree
	tree := &MerkleTree{
		Root:           root,
		MerkleRootHash: t.Nodes[root].Hash,
		Leafs:          t.Leafs,
		Nodes:          t.Nodes,
		HashStrategy:   t.HashStrategy,
	}
	if err != nil {
		return nil
	}

	// Return the tree
	return tree
}

// updateLeafsSlice add a new leaf to data structure, if odd number of leafs are present duplicates the last leaf
func updateLeafsSlice(leafs []int, nodes []*Node, leaf *Node) ([]int, []*Node) {
	if len(leafs) > 0 && nodes[leafs[len(leafs)-1]].Dup {
		leafs = leafs[:len(leafs)-1]
	}
	var newNodes []*Node
	for i := 0; i < len(leafs); i++ {
		newNodes = append(newNodes, nodes[leafs[i]])
	}

	leafs = append(leafs, len(leafs))
	newNodes = append(newNodes, leaf)

	// If odd number of leafs are present, duplicates the last leaf
	if len(leafs)%2 != 0 {
		lastLeaf := newNodes[leafs[len(leafs)-1]]
		newLeaf := &Node{
			Parent: -1,
			Left:   -1,
			Right:  -1,
			Leaf:   true,
			Dup:    true,
			Hash:   lastLeaf.Hash,
		}
		newLeafIndex := len(newNodes)
		newNodes = append(newNodes, newLeaf)
		leafs = append(leafs, newLeafIndex)
	}
	return leafs, newNodes
}

// buildIntermediate recursively builds intermediate nodes until the root node is reached
func buildIntermediate(nodesSlice []int, t *MerkleTree) (int, error) {

	// If node slice have odd length make it even to loop through
	if len(nodesSlice)%2 == 1 {
		nodesSlice = append(nodesSlice, nodesSlice[len(nodesSlice)-1])
	}
	h := sha256.New()
	var parents []int
	for i := 0; i < len(nodesSlice); i += 2 {
		left := nodesSlice[i]
		right := nodesSlice[i+1]
		h.Reset()
		h.Write(t.Nodes[left].Hash)
		h.Write(t.Nodes[right].Hash)
		parent := &Node{
			Parent: -1,
			Left:   left,
			Right:  right,
			Leaf:   false,
			Dup:    false,
			Hash:   h.Sum(nil),
		}
		parentIndex := len(t.Nodes)
		t.Nodes = append(t.Nodes, parent)
		t.Nodes[left].Parent = parentIndex
		t.Nodes[right].Parent = parentIndex
		parents = append(parents, parentIndex)
	}
	if len(parents) == 1 {
		return parents[0], nil
	} else {
		return buildIntermediate(parents, t)
	}
}

// VerifyContentHash verifies the hash of a given content against the Merkle tree
func (t *MerkleTree) VerifyContentHash(content []byte, rootHash []byte) bool {
	// Find the leaf node that contains the matching hash
	leafNodeIndex := -1
	hashVal := t.hashByteSlice(content)
	for _, leafIndex := range t.Leafs {
		if bytes.Equal(t.Nodes[leafIndex].Hash, hashVal) {
			leafNodeIndex = leafIndex
			break
		}
	}
	if leafNodeIndex == -1 {
		return false
	}

	// Traverse the tree from the leaf node up to the root node
	nodeIndex := leafNodeIndex
	for t.Nodes[nodeIndex].Parent != -1 {
		parentIndex := t.Nodes[nodeIndex].Parent
		leftIndex := t.Nodes[parentIndex].Left
		rightIndex := t.Nodes[parentIndex].Right

		// Hash the hashes of the left and right child nodes
		h := sha256.New()
		h.Write(t.Nodes[leftIndex].Hash)
		h.Write(t.Nodes[rightIndex].Hash)
		hashValue := h.Sum(nil)[:]

		// If the computed hash matches the parent node's hash, continue up the tree
		if bytes.Equal(hashValue, t.Nodes[parentIndex].Hash) {
			nodeIndex = parentIndex
		} else {
			// If the computed hash does not match the parent node's hash, the content has been tampered with
			return false
		}
	}

	// If we've reached the root node and its hash matches the expected root hash, the content has not been tampered with
	return bytes.Equal(t.Nodes[nodeIndex].Hash, rootHash)
}

// PrintTreeNodes returns a string representation of the Merkle tree
func (t *MerkleTree) PrintTreeNodes() {
	s := ""
	for _, n := range t.Nodes {
		s += fmt.Sprintf("%v\n", n)
	}
	fmt.Println(s)
}
