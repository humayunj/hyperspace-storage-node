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
	Root       []byte
	LeafIndex  int
	Proof      [][]byte
	Data       []byte
	Directions []uint32
}

func ensureEven(hashes [][]byte) [][]byte {
	if len(hashes)%2 != 0 {
		a := make([]byte, 0)
		copy(a, hashes[len(hashes)-1])
		hashes = append(hashes, a)
	}
	return hashes
}

func _generateLevel(hashes [][]byte, tree [][][]byte) [][][]byte {
	if len(hashes) == 1 {
		return nil
	}
	hashes = ensureEven(hashes)

	combinedHashes := make([][]byte, 0)

	for i := 0; i < len(hashes); i += 2 {
		hashConcat := bytes.Join([][]byte{hashes[i], hashes[i+1]}, []byte{})
		hash := keccak256.New().Hash(hashConcat)
		combinedHashes = append(combinedHashes, hash)
	}
	tree = append(tree, combinedHashes)
	return _generateLevel(combinedHashes, tree)
}

func GenerateMerkleTree(leaves [][]byte) [][][]byte {
	if len(leaves) == 0 {
		return nil
	}
	tree := make([][][]byte, 0)

	tree = append(tree, leaves)

	tree = _generateLevel(leaves, tree)
	return tree
}

func GenerateMerkleProof(leafIndex uint, leaves [][]byte) ([][]byte, []uint32) {
	if len(leaves) == 0 {
		return nil, nil
	}
	tree := GenerateMerkleTree(leaves)
	merkleProof := make([][]byte, 0)
	directions := make([]uint32, 0)
	merkleProof = append(merkleProof, leaves[leafIndex])
	d := uint32(0)
	if leafIndex%2 != 0 {
		d = 1
	}
	directions = append(directions, d)

	hashIndex := leafIndex
	for level := 0; level < len(tree)-1; level++ {
		isLeftChild := hashIndex%2 == 0

		siblingDirection := uint32(0)
		if isLeftChild {
			siblingDirection = 1
		}
		siblingIndex := hashIndex - 1
		if isLeftChild {
			siblingIndex = hashIndex + 1
		}
		//   const siblingNode = {
		// 	hash: tree[level][siblingIndex],
		// 	direction: siblingDirection,
		//   };
		merkleProof = append(merkleProof, tree[level][siblingIndex])
		directions = append(directions, siblingDirection)
		//   merkleProof.push(siblingNode);
		hashIndex = (hashIndex / 2)
	}
	return merkleProof, directions
}
func GetMerkleProof(leaves [][]byte, leaveInd int) ([][]byte, []uint32, error) {
	if leaveInd >= len(leaves) {
		return nil, nil, errors.New("leave index out of bounds")
	}

	path, dir := GenerateMerkleProof(uint(leaveInd), leaves)
	return path, dir, nil
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
func createLeavesExt(file string, segmentsCount uint32, segmentIndex int) ([][]byte, []byte, error) {

	stats, err := os.Stat(file)
	if err != nil {
		return nil, nil, err
	}
	fileSize := stats.Size()
	o_segs := segmentsCount
	segmentsCount = uint32(math.Ceil(float64(fileSize) / 1024))

	if o_segs != segmentsCount {
		printLn(">Seg count mismatch ", o_segs, "!=", segmentsCount)
	} else {
		printLn("Seg count match:", segmentsCount)
	}

	lastChunkSize := (fileSize % 1024)

	var segmentSize uint32

	if segmentsCount == 1 {
		segmentSize = uint32(lastChunkSize)
	} else {
		segmentSize = uint32(math.Floor(float64(fileSize-lastChunkSize) / float64(segmentsCount-1)))
	}
	var segments [][]byte
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
			segments = append(segments, keccak256.New().Hash(c))
			if i == (segmentIndex) {
				segData = *bytes.NewBuffer(c)
			}
			i += 1
		}
	}

	return segments, segData.Bytes(), nil

}
func createLeaves(file string, segmentsCount uint32) ([][]byte, error) {
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

	proof, directions, err := GetMerkleProof(leaves, segmentIndex)
	if err != nil {
		return nil, err
	}
	merkleProof = new(MerkleProof)
	merkleProof.LeafIndex = segmentIndex
	merkleProof.Proof = proof
	merkleProof.Data = segmentData
	merkleProof.Root = tree.Root.Hash
	merkleProof.Directions = directions

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
