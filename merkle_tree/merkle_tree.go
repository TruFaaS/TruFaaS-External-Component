package merkle_tree

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"
)

// MerkleTree represents a Merkle tree
type MerkleTree struct {
	Nodes          []*Node // Nodes in the Merkle tree
	RootIndex      int     // RootIndex is the index of the root node in the nodes slice
	LeafCount      int     // LeafCount holds the number of leafs
	MerkleRootHash []byte  // MerkleRootHash is the hash of the Merkle tree root
}

// Node represents a node in the Merkle tree
type Node struct {
	Parent int    // Index of the parent node in the nodes slice, -1 if nil
	Left   int    // Index of the left child node in the nodes slice, -1 if nil
	Right  int    // Index of the right child node in the nodes slice, -1 if nil
	Leaf   bool   // Indicates whether the node is a leaf node
	Dup    bool   // Indicates whether the node is a duplicate leaf node
	Hash   []byte // Hash value of the node
}

// NewTree creates a new Merkle Tree
func NewTree() *MerkleTree {
	t := &MerkleTree{Nodes: make([]*Node, 0)}
	return t
}

// GetMerkleRoot returns the root hash of the Merkle tree
func (t *MerkleTree) GetMerkleRoot() []byte {
	return t.MerkleRootHash
}

// hashByteSlice returns the hash of a byte slice using the hash function
func (t *MerkleTree) hashByteSlice(data []byte) []byte {
	h := NewHashFunc()
	h.Write(data)
	return h.Sum(nil)[:]
}

// AppendNewContent builds a new tree with the new content and return the tree
func (t *MerkleTree) AppendNewContent(content []byte) *MerkleTree {

	// Create new leaf
	leaf := &Node{
		Parent: -1, // Parent not set
		Left:   -1, // Left not set
		Right:  -1, // Right not set
		Leaf:   true,
		Dup:    false,
		Hash:   t.hashByteSlice(content), // Hash value by hashing the content
	}

	// Update the leaf count and nodes of tree
	t.LeafCount, t.Nodes = updateLeafsAndNodes(t.LeafCount, t.Nodes, leaf)

	// Leaf indices from [0,1..,t.leafCount]
	leafIndices := make([]int, t.LeafCount)
	for i := range leafIndices {
		leafIndices[i] = i
	}
	// Recursively build intermediate nodes until the root node is reached
	t.RootIndex, _ = buildIntermediate(leafIndices, t)
	t.MerkleRootHash = t.Nodes[t.RootIndex].Hash

	// Return the tree
	return t
}

// updateLeafsAndNodes add a new leaf to data structure and clear previously build intermediate nodes
func updateLeafsAndNodes(leafCount int, nodes []*Node, leaf *Node) (int, []*Node) {

	// This list stores list of leaf objects
	var newLeafNodes []*Node

	// If there are leafs check if last is a duplicate if so remove it
	if leafCount > 0 && nodes[leafCount-1].Dup {
		leafCount = leafCount - 1
	}

	// Add all the leafs to newLeafNodes from the previous nodes list
	for i := 0; i < leafCount; i++ {
		newLeafNodes = append(newLeafNodes, nodes[i])
	}

	// Add the new leaf as a leaf
	newLeafNodes = append(newLeafNodes, leaf)
	leafCount += 1

	// If odd number of leafIndices are present, duplicates the last leaf
	if leafCount%2 != 0 {
		lastLeaf := newLeafNodes[leafCount-1]
		newLeaf := &Node{
			Parent: -1,
			Left:   -1,
			Right:  -1,
			Leaf:   true,
			Dup:    true,          // Duplicate
			Hash:   lastLeaf.Hash, // Same hash as it is a duplicate
		}
		newLeafNodes = append(newLeafNodes, newLeaf)
		leafCount += 1
	}
	return leafCount, newLeafNodes
}

// buildIntermediate recursively builds intermediate nodes until the root node is reached
func buildIntermediate(nodesIndexSlice []int, t *MerkleTree) (int, error) {

	// Hash start
	h := NewHashFunc()

	// If node slice have odd length make it even to loop through
	if len(nodesIndexSlice)%2 == 1 {
		nodesIndexSlice = append(nodesIndexSlice, nodesIndexSlice[len(nodesIndexSlice)-1])
	}

	var parentIndices []int

	// Attempt to create parent for pairs of leafs
	for i := 0; i < len(nodesIndexSlice); i += 2 {
		left := nodesIndexSlice[i]    // Get the left child index
		right := nodesIndexSlice[i+1] // Get the right child index

		h.Reset()
		h.Write(t.Nodes[left].Hash)  // Get the hash of left child
		h.Write(t.Nodes[right].Hash) // Get the hash of right child
		parent := &Node{
			Parent: -1,
			Left:   left,
			Right:  right,
			Leaf:   false,
			Dup:    false,
			Hash:   h.Sum(nil), // Hash of hashes
		}
		parentIndex := len(t.Nodes)
		t.Nodes = append(t.Nodes, parent)                  // append the created node to tree
		t.Nodes[left].Parent = parentIndex                 // update the parent of left child
		t.Nodes[right].Parent = parentIndex                // update the parent of right child
		parentIndices = append(parentIndices, parentIndex) // add the parent index
	}
	if len(parentIndices) == 1 { // Root is reached
		return parentIndices[0], nil // return root index
	} else {
		return buildIntermediate(parentIndices, t) // build the next level
	}
}

// VerifyContentHash verifies the hash of a given content against the Merkle tree
func (t *MerkleTree) VerifyContentHash(content []byte, rootHash []byte) bool {

	// Find the leaf node that contains the matching hash
	leafNodeIndex := -1                 // No index set
	hashVal := t.hashByteSlice(content) // Hash the given content
	for i := 0; i < t.LeafCount; i++ {
		if bytes.Equal(t.Nodes[i].Hash, hashVal) { // If current hash matches any of leaf hashes
			leafNodeIndex = i
			break
		}
	}
	if leafNodeIndex == -1 { // If no leaf found the content verification fails
		return false
	}

	// Traverse the tree from the leaf node up to the root node
	nodeIndex := leafNodeIndex
	for t.Nodes[nodeIndex].Parent != -1 { // Till root is reached
		parentIndex := t.Nodes[nodeIndex].Parent
		leftIndex := t.Nodes[parentIndex].Left
		rightIndex := t.Nodes[parentIndex].Right

		// Hash the hashes of the left and right child nodes
		h := NewHashFunc()
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
	// We ultimately check the root hash with the passed hash to this method
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

// NewHashFunc Returns the hash function, this is the only place to be changed to change the hash func
func NewHashFunc() hash.Hash {
	return sha256.New()
}
