package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"math"
	"os"

	"github.com/cbergoon/merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"
)

type MerkleProof struct {
	Root      []byte
	LeafIndex int
	Proof     [][]byte
	Data      []byte
}

func GetMerkleProof(m *merkletree.MerkleTree, leaveInd int) ([][]byte, error) {
	if leaveInd >= len(m.Leafs) {
		return nil, errors.New("leave index out of bounds")
	}
	current := m.Leafs[leaveInd]

	currentParent := current.Parent
	var merklePath [][]byte
	// var index []int64
	for currentParent != nil {
		if bytes.Equal(currentParent.Left.Hash, current.Hash) {
			merklePath = append(merklePath, currentParent.Right.Hash)
			// index = append(index, 1) // right leaf
		} else {
			merklePath = append(merklePath, currentParent.Left.Hash)
			// index = append(index, 0) // left leaf
		}
		current = currentParent
		currentParent = currentParent.Parent
	}
	return merklePath, nil

}

type TContent struct {
	x []byte
}

func (t TContent) CalculateHash() ([]byte, error) {
	// printLn(">x:", hex.EncodeToString(t.x))

	// printLn("Hash i: ", hex.EncodeToString(t.x))
	// return keccak256.New().Hash(t.x), nil
	return t.x, nil
}

func (t TContent) Equals(other merkletree.Content) (bool, error) {
	return bytes.Equal(t.x, other.(TContent).x), nil
}
func createLeavesExt(file string, segmentsCount uint32, segmentIndex int) ([]merkletree.Content, []byte, error) {

	stats, err := os.Stat(file)
	if err != nil {
		return nil, nil, err
	}
	fileSize := stats.Size()
	segmentsCount = uint32(math.Ceil(float64(fileSize) / 1024))

	lastChunkSize := (fileSize % 1024)

	var segmentSize uint32

	if segmentsCount == 1 {
		segmentSize = uint32(lastChunkSize)
	} else {
		segmentSize = uint32(math.Floor(float64(fileSize-lastChunkSize) / float64(segmentsCount-1)))
	}
	var segments []merkletree.Content
	f, err := os.Open(file)
	if err != nil {
		return nil, nil, err
	}
	var segData bytes.Buffer
	i := 0
	for {
		chunk := make([]byte, segmentSize)
		n, err := f.Read(chunk)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, err
		}
		if n > 0 {
			// printLn(">seg:", hex.EncodeToString(chunk[:n])[:30])
			c := chunk[:n]
			segments = append(segments, TContent{keccak256.New().Hash(c)})
			if i == (segmentIndex) {
				segData = *bytes.NewBuffer(c)
			}
			i += 1
		}
	}

	return segments, segData.Bytes(), nil

}
func createLeaves(file string, segmentsCount uint32) ([]merkletree.Content, error) {
	leaves, _, err := createLeavesExt(file, segmentsCount, -1)
	return leaves, err
}

type Keccak256Hash struct {
	buffer []byte
}

func (k *Keccak256Hash) Write(p []byte) (n int, err error) {
	k.buffer = append(k.buffer, p...)
	return len(p), nil
}
func (k *Keccak256Hash) Sum(b []byte) []byte {
	if b != nil {
		k.Write(b)
	}
	// printLn("Hash: ", hex.EncodeToString(keccak256.New().Hash(k.buffer)))

	return keccak256.New().Hash(k.buffer)
}
func (k *Keccak256Hash) Reset() {
	k.buffer = make([]byte, 0)
}
func (k *Keccak256Hash) Size() int {
	return 32
}

func (k *Keccak256Hash) BlockSize() int {
	return len(k.buffer)
}

func ComputeMerkleProof(filePath string, segments uint32, segmentIndex int) (merkleProof *MerkleProof, err error) {

	leaves, segmentData, err := createLeavesExt(filePath, uint32(segments), segmentIndex)

	if err != nil {
		return nil, err
	}

	k := new(Keccak256Hash)
	k.buffer = make([]byte, 0)

	tree, err := merkletree.NewTreeWithHashStrategy(leaves, func() hash.Hash {
		return &Keccak256Hash{buffer: make([]byte, 0)}
	})

	if err != nil {
		return nil, err
	}
	for _, t := range tree.Leafs {
		printLn(hex.EncodeToString(t.Hash))
	}
	proof, err := GetMerkleProof(tree, segmentIndex)
	if err != nil {
		return nil, err
	}
	merkleProof = new(MerkleProof)
	merkleProof.LeafIndex = segmentIndex
	merkleProof.Proof = proof
	merkleProof.Data = segmentData
	merkleProof.Root = tree.Root.Hash

	return merkleProof, err
}

func ComputeFileRootMerkle(filePath string, segments uint32) (root string, err error) {

	leaves, err := createLeaves(filePath, uint32(segments))

	if err != nil {
		return "", err
	}

	k := new(Keccak256Hash)
	k.buffer = make([]byte, 0)

	tree, err := merkletree.NewTreeWithHashStrategy(leaves, func() hash.Hash {
		return &Keccak256Hash{buffer: make([]byte, 0)}
	})

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(tree.MerkleRoot()), nil

}
