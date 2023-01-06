package main

import (
	"bytes"
	"encoding/hex"
	"hash"
	"io"
	"math"
	"os"

	"github.com/cbergoon/merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"
)

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

func createLeaves(file string, segmentsCount uint32) ([]merkletree.Content, error) {

	stats, err := os.Stat(file)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	for {
		chunk := make([]byte, segmentSize)
		n, err := f.Read(chunk)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if n > 0 {
			// printLn(">seg:", hex.EncodeToString(chunk[:n])[:30])

			segments = append(segments, TContent{keccak256.New().Hash(chunk[:n])})
		}
	}

	return segments, nil

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
