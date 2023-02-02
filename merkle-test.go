package main

import (
	"encoding/hex"

	"github.com/wealdtech/go-merkletree"
	mt "github.com/wealdtech/go-merkletree"
)

// Example using the Merkle tree to generate and verify proofs.
func mainMerkle() {
	// Data for the tree
	data := [][]byte{
		[]byte("Foo"),
		[]byte("Bar"),
		[]byte("Baz"),
		[]byte("Baz"),
		[]byte("Baz"),
		[]byte("Baz"),
	}

	// Create the tree
	tree, err := mt.NewTree(mt.WithData(data))
	if err != nil {
		panic(err)
	}

	// Fetch the root hash of the tree
	root := tree.Root()

	baz := data[2]
	// Generate a proof for 'Baz'
	proof, err := tree.GenerateProof(baz, 0)
	if err != nil {
		panic(err)
	}

	for _, i := range proof.Hashes {
		println(hex.EncodeToString(i))
	}
	// Verify the proof for 'Baz'
	verified, err := merkletree.VerifyProof(baz, false, proof, [][]byte{root})
	if err != nil {
		panic(err)
	}
	if !verified {
		panic("failed to verify proof for Baz")
	}
}
