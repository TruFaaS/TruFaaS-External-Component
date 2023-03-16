package merkle_tree

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"
)

// MerkleTree represents a Merkle tree
type MerkleTree struct {
	Root           *Node            // Root node of the Merkle tree
	merkleRootHash []byte           // Root hash of the Merkle tree
	Leafs          []*Node          // Leaf nodes of the Merkle tree
	hashStrategy   func() hash.Hash // Hash function used to compute Merkle tree hashes
}

// Node represents a node in the Merkle tree
type Node struct {
	Parent *Node  // Pointer to the parent node
	Left   *Node  // Pointer to the left child node
	Right  *Node  // Pointer to the right child node
	leaf   bool   // Indicates whether the node is a leaf node
	dup    bool   // Indicates whether the node is a duplicate leaf node
	Hash   []byte // Hash value of the node
}

// NewTree creates a new Merkle Tree with default hash strategy of SHA256
func NewTree() *MerkleTree {
	var defaultHashStrategy = sha256.New
	t := &MerkleTree{hashStrategy: defaultHashStrategy, Leafs: make([]*Node, 0)}
	return t
}

// MerkleRoot returns the root hash of the Merkle tree
func (t *MerkleTree) MerkleRoot() []byte {
	return t.merkleRootHash
}

// hashByteSlice returns the hash of a byte slice using the hash function specified in the Merkle tree
func (t *MerkleTree) hashByteSlice(data []byte) []byte {
	h := t.hashStrategy()
	h.Write(data)
	return h.Sum(nil)
}

// AppendNewContent builds a new tree with the new content and return the tree
func (t *MerkleTree) AppendNewContent(content []byte) *MerkleTree {

	leaf := &Node{
		Parent: nil,
		Left:   nil,
		Right:  nil,
		leaf:   true,
		dup:    false,
		Hash:   t.hashByteSlice(content),
	}
	leafs := updateLeafsSlice(t.Leafs, leaf)

	// Recursively build intermediate nodes until the root node is reached
	root, err := buildIntermediate(leafs, t)

	// Create the new tree
	tree := &MerkleTree{
		Root:           root,
		merkleRootHash: root.Hash,
		Leafs:          leafs,
		hashStrategy:   t.hashStrategy,
	}
	if err != nil {
		return nil
	}

	// Return the tree
	return tree
}

// updateLeafsSlice add a new leaf to data structure, if odd number of leafs are present duplicates the last leaf
func updateLeafsSlice(leafs []*Node, leaf *Node) []*Node {
	if len(leafs) > 0 && leafs[len(leafs)-1].dup {
		leafs = leafs[:len(leafs)-1]
	}
	leafs = append(leafs, leaf)
	if len(leafs)%2 == 1 {
		duplicate := &Node{
			Hash: leafs[len(leafs)-1].Hash,
			leaf: true,
			dup:  true,
		}
		leafs = append(leafs, duplicate)
	}
	return leafs
}

// buildIntermediate constructs intermediate nodes in the Merkle tree
func buildIntermediate(nodes []*Node, t *MerkleTree) (*Node, error) {
	h := t.hashStrategy()
	var newIntermediateNodes []*Node
	for i := 0; i < len(nodes); i += 2 {
		var left, right = i, i + 1
		if i+1 == len(nodes) {
			right = i
		}
		h.Reset()
		h.Write(nodes[left].Hash)
		h.Write(nodes[right].Hash)
		n := &Node{
			Left:  nodes[left],
			Right: nodes[right],
			Hash:  h.Sum(nil),
		}
		newIntermediateNodes = append(newIntermediateNodes, n)
		nodes[left].Parent = n
		nodes[right].Parent = n
		if len(nodes) == 2 {
			return n, nil
		}
	}
	return buildIntermediate(newIntermediateNodes, t)
}

// VerifyContentHash verifies the hash of a given content against the Merkle tree
func (t *MerkleTree) VerifyContentHash(content []byte, rootHash []byte) bool {
	// Find the leaf node that contains the matching hash
	var node *Node
	for _, n := range t.Leafs {
		if bytes.Equal(n.Hash, t.hashByteSlice(content)) {
			node = n
		}
	}
	if node == nil {
		return false
	}

	// Traverse the tree from the leaf node up to the root node
	for node.Parent != nil {
		left := node.Parent.Left
		right := node.Parent.Right

		// Hash the hashes of the left and right chil d nodes
		h := t.hashStrategy()
		h.Write(left.Hash)
		h.Write(right.Hash)
		hashVal := h.Sum(nil)

		// If the computed hash matches the parent node's hash, continue up the tree
		if bytes.Equal(hashVal, node.Parent.Hash) {
			node = node.Parent
		} else {
			// If the computed hash does not match the parent node's hash, the content has been tampered with
			return false
		}
	}

	// If we've reached the root node and its hash matches the expected root hash, the content has not been tampered with
	return bytes.Equal(node.Hash, rootHash)
}

// PrintLeafs prints leafs of the tree
func (t *MerkleTree) PrintLeafs() {
	s := ""
	for _, l := range t.Leafs {
		s += fmt.Sprint(l)
		s += "\n"
	}
	fmt.Println(s)
}