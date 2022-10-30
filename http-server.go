package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	color "github.com/fatih/color"

	"io"
	"log"
	"net/http"
)

func processUpload(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		http.Error(w, ("Provide Authorization token"), http.StatusInternalServerError)
		return
	}

	authHeader = authHeader[7:] // Skip "bearer "
	if len(authHeader) == 0 {
		http.Error(w, ("Provide valid JWT token"), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Authorization: %v\n", authHeader)

	r.ParseMultipartForm(1024 * 1024 * 10) // 10 mb

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer file.Close()
	// fmt.Printf("File Size: %v\n", header.Size)

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	const key = "3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"
	hash := h.Sum(nil)
	err = FS.AddFile(file, hash, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Print(err)
		return
	}

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
	fmt.Fprintln(tw, "Key\t"+key)
	fmt.Fprintln(tw, "Size\t", header.Size)
	fmt.Fprintln(tw, "Hash\t", hex.EncodeToString(hash))
	tw.Flush()
	color.Unset()
	log.Println("TODO: Conclude Transaction with smart contract")

}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		processUpload(w, r)
	default:
		w.Write([]byte("METHOD NOT ALLOWED"))
	}
}

func runHTTPServer() {
	color.Set(color.FgBlue)

	log.Println("HTTP server listening on port 5000")
	color.Unset()
	http.HandleFunc("/upload", handler)
	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		log.Fatal(err)
	}
}
