package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"math"
	"os"

	"github.com/cbergoon/merkletree"
	gm "github.com/wealdtech/go-merkletree/keccak256"
)

func keccak256(b []byte) []byte {
	return gm.New().Hash(b)
}

type MerkleProof struct {
	Root       []byte
	LeafIndex  int
	Proof      [][]byte
	Data       []byte
	Directions []uint32
}

func ensureEven(hashes [][]byte) [][]byte {
	if len(hashes)%2 != 0 {
		// a := make([]byte, 1)
		// copy(a, hashes[len(hashes)-1])
		printLn("Appending as ", len(hashes))
		h := hashes[:]
		h = append(h, (hashes)[len(hashes)-1])
		return h
	}

	return hashes
}

func GenerateMerkleRoot(leaves [][]byte) []byte {
	if len(leaves) == 0 {
		return nil
	}
	printLn("Pre Leaves length: ", len(leaves))

	leaves = ensureEven(leaves)
	printLn("Leaves length: ", len(leaves))
	// printLn("Leaves", leaves)
	var combinedHashes [][]byte
	for i := 0; i < len(leaves); i += 2 {
		// printLn(hex.EncodeToString(leaves[i]), " ", hex.EncodeToString(leaves[i+1]))
		hashConcat := bytes.Join([][]byte{leaves[i], leaves[i+1]}, []byte{})
		hash := keccak256(hashConcat)
		// printLn(">Hash", hash)
		combinedHashes = append(combinedHashes, hash)
	}
	if len(combinedHashes) == 1 {
		return combinedHashes[0]
	}
	return GenerateMerkleRoot(combinedHashes)
}

func _generateLevel(hashes [][]byte, tree [][][]byte) [][][]byte {
	if len(hashes) == 1 {
		printLn("return 1")

		printLn(">>>>")
		for i, r := range tree {
			printLn("Level: ", i, " len:", len(r))
			if len(r) < 15 {
				for _, n := range r {
					printLn(n)
				}
				printLn(("\n"))
			}
		}

		return tree[:]
	}
	printLn("Pre hash length: ", len(hashes))
	hashes = ensureEven(hashes)

	printLn("Leaves length: ", len(hashes))

	combinedHashes := make([][]byte, 0)

	for i := 0; i < len(hashes); i += 2 {
		hashConcat := bytes.Join([][]byte{hashes[i], hashes[i+1]}, []byte{})
		hash := keccak256(hashConcat)
		combinedHashes = append(combinedHashes, hash[:])
	}
	printLn(">Combined hashes length:", len(combinedHashes))
	cHashes := combinedHashes[:]
	treeN := tree[:]
	treeN = append(treeN, cHashes)
	return _generateLevel(cHashes, treeN)
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
	for i, r := range tree {
		printLn("Level: ", i, " len:", len(r))
		if len(r) < 15 {
			for _, n := range r {
				printLn(n)
			}
			printLn(("\n"))
		}
	}
	merkleProof := make([][]byte, 0, 10)
	directions := make([]uint32, 0, 10)
	merkleProof = append(merkleProof, leaves[leafIndex])
	d := uint32(0)
	if leafIndex%2 != 0 {
		d = 1
	}
	directions = append(directions, d)

	hashIndex := leafIndex
	printLn("Tree Length:", len(tree))
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

		printLn("Tree([Level]length)", len(tree[level]))
		if int(siblingIndex) >= len(tree[level]) {
			printLn("sibling out of range: ", siblingIndex, len(tree[level]))
			return nil, nil
		}

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
			segments = append(segments, keccak256(c))
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

func ComputeMerkleProof(filePath string, segments uint32, segmentIndex int) (merkleProof *MerkleProof, err error) {

	leaves, segmentData, err := createLeavesExt(filePath, uint32(segments), segmentIndex)
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
	merkleProof.Root = GenerateMerkleRoot(leaves)
	merkleProof.Directions = directions

	return merkleProof, err
}

func ComputeFileRootMerkle(filePath string, segments uint32) (string, error) {

	leaves, err := createLeaves(filePath, uint32(segments))

	if err != nil {
		return "", err
	}

	root := GenerateMerkleRoot(leaves)

	return hex.EncodeToString(root), nil

}
