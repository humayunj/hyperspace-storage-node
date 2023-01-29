package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/ethereum/go-ethereum/crypto"
	color "github.com/fatih/color"

	"io"
	"log"
	"net/http"
)

func processUpload(w http.ResponseWriter, r *http.Request) {
	printLn("Upload Begin")
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		http.Error(w, ("Provide Authorization token"), http.StatusBadRequest)
		return
	}

	authHeader = authHeader[7:] // Skip "bearer "
	if len(authHeader) == 0 {
		http.Error(w, ("Provide valid JWT token"), http.StatusBadRequest)
		return
	}
	// fmt.Printf("Authorization: %v\n", authHeader)
	claims, ok := JFS.ParseToken(authHeader)
	if ok != nil {
		http.Error(w, "JWT token is not valid", http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(1024 * 1024 * 1024 * 10) // 10 GB

	file, header, err := r.FormFile("file")
	if err != nil {
		printLn("'file' key form data read error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer file.Close()
	// fmt.Printf("File Size: %v\n", header.Size)

	// h := sha256.New()

	dir := "uploads/"
	if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	printLn("User address hex", claims.UserAddress[2:])
	addressBytes, err := hex.DecodeString(claims.UserAddress[2:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	hashBytes, err := hex.DecodeString(claims.FileHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	keyBytes := crypto.Keccak256(bytes.Join([][]byte{addressBytes, hashBytes}, []byte{}))

	keyHex := hex.EncodeToString(keyBytes)

	printLn("Key Hex: ", keyHex)

	path := strings.Join([]string{"uploads/", keyHex}, "")

	newFile, err := os.Create(path)
	if err != nil {
		log.Print(err)
		// return err
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	printLn("Copy: ", newFile, " => ", file)

	if _, err := io.Copy(newFile, file); err != nil {
		log.Print(err)
		// return err
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = file.Close()
	if err != nil {
		log.Println("Failed to close file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	printLn("Computing merkle root hash...")
	root, err := ComputeFileRootMerkle(path, uint32(claims.SegmentsCount))
	if err != nil {
		printLn("Failed to compute merkle root hash")
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if root != claims.FileHash {
		printLn(fmt.Sprintf("Computed merkle root hash (%v) mismatch claimed file hash (%v)", root, claims.FileHash))
		http.Error(w, "claimed hash mismatch computed hash", http.StatusInternalServerError)
		return
	}

	bid := new(big.Int)
	bid, success := bid.SetString(claims.Bid, 10)
	if !success {
		log.Println("Failed to parse bid")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = CG.ConcludeTransaction(claims.UserAddress,
		claims.FileHash, uint32(claims.FileSize),
		claims.TimeStart,
		claims.TimeEnd,
		claims.ProveTimeout,
		claims.ConcludeTimeout,
		claims.SegmentsCount, bid)
	if err != nil {
		log.Println(err.Error())

		http.Error(w, "failed to conclude tx", http.StatusInternalServerError)
		return
	}

	err = DBS.InsertTransaction(TransactionParams{
		FileKey:            keyHex,
		UserAddress:        claims.UserAddress,
		FileMerkleRootHash: claims.FileHash,
		FileName:           "RESERVED",
		FileSize:           claims.FileSize,
		Status:             TRANSACTION_STATUS_PENDING,
		BidPrice:           claims.Bid,
		UploadedAt:         claims.TimeStart,
		ExpiresAt:          claims.TimeEnd,
	})

	if err != nil {
		color.Set(color.FgRed)
		printLn(err)
		color.Unset()
		os.Remove(path)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// if _, err := io.Copy(h, file); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	log.Print(err)
	// 	return
	// }
	// const key = "3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"
	// hash := h.Sum(nil)
	// err = FS.AddFile(file, hash, key)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	log.Print(err)
	// 	return
	// }

	// CG.ConcludeTransaction(claims.UserAddress,claims)

	printLn("Resp")
	type TResp struct {
		Ok bool `json:"ok"`
	}
	res := TResp{
		Ok: true,
	}
	re, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	fmt.Fprint(w, string(re))

	// Call Conclude Transaction
	log.Printf("Uploaded new file")
	color.Set(color.FgGreen)
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	// fmt.Fprintln(tw, "Key\t"+key)
	fmt.Fprintln(tw, "Size\t", header.Size)
	// fmt.Fprintln(tw, "Hash\t", hex.EncodeToString(hash))
	tw.Flush()
	color.Unset()
	log.Println("TODO: Finish Transaction with smart contract")

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		processUpload(w, r)
	default:
		w.Write([]byte("METHOD NOT ALLOWED"))
	}
}

func processDownload(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	printLn("Path: ", p)
	fileKey := p[5:]
	printLn("FileKey:", fileKey)

	if len(fileKey) == 0 {
		http.Error(w, "filekey is invalid", http.StatusBadRequest)
		return
	}
	tx, err := DBS.GetTransaction(fileKey)

	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	w.Header().Set("X-MERKLE-ROOT-HASH", tx.FileMerkleRootHash)

	http.ServeFile(w, r, "./uploads/"+tx.FileKey)

}
func getHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		processDownload(w, r)
	default:
		w.Write([]byte("METHOD NOT ALLOWED"))
	}
}
func runHTTPServer() {
	color.Set(color.FgMagenta)

	log.Println("Listening HTTP on 5555")
	color.Unset()
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/get/", getHandler)
	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Fatal(err)
	}
}
