package main

import (
	"fmt"
	"log"
	"os"

	"text/tabwriter"

	color "github.com/fatih/color"
	"google.golang.org/grpc"
)

var FS FileService
var JFS JWTFileService
var CG *ContractGateway
var NC *NodeConfig
var GRPCServer *grpc.Server
var DBS *DBService

func initContract(contractAddress string) {
	color.Set(color.FgYellow)
	CG, cgErr := NewContractRPC(contractAddress)
	if cgErr != nil {
		log.Fatalln(cgErr)
	}

	log.Println("Fetching Contract Stats...")
	color.Set(color.FgGreen)
	contractStats, err := CG.GetContractStats()
	if err != nil {
		color.Set(color.FgRed)
		log.Println("Failed to fetch contract stats")
		log.Println(err)
		color.Unset()
	} else {

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintln(tw, "Balance\t", weiToEther(contractStats.Balance))
		fmt.Fprintln(tw, "Locked Collateral\t", contractStats.LockedCollateral)
		fmt.Fprintln(tw, "Mappings Length\t", contractStats.MappingLength)

		tw.Flush()

	}
	color.Unset()
}

func openDB() {

	println("Opening DB")

	if _, err := DBS.connect(); err != nil {
		panic(err)
	}
}

func printLn(v ...interface{}) {
	log.Print(v...)
}

func tokenTest() {
	token, err := JFS.CreateFileToken(FileTokenParams{
		FileHash: "demoFileHash",
	})

	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println(token)

	params, err := JFS.ParseToken(token)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println("Hash", params.FileHash)

}
func main() {
	// To Get the token
	// tokenTest()
	// return
	println("Starting...")
	openDB()
	initContract("")

	println("Loading Configuration")

	NC = LoadConfig()
	color.Set(color.FgGreen)
	println("Total Storage:", NC.TotalStorage, "B")
	color.Unset()

	// FS.Open()
	// println("Fetching Already Uploaded Files")
	// f, err := FS.GetAllFiles()
	// if err != nil {
	// 	panic(err)
	// }

	// tw := tabwriter.NewWriter(priter(), 0, 0, 1, ' ', 0)

	// color.Set(color.FgGreen)
	// for i, file := range f {

	// 	fmt.Fprintln(tw, fmt.Sprint(i)+":\t")
	// 	fmt.Fprintln(tw, "Path:\t"+file.Path)
	// 	fmt.Fprintln(tw, "Hash:\t"+hex.EncodeToString(file.Hash))
	// 	fmt.Fprintln(tw, "Added At:\t"+fmt.Sprint(file.AddedAt))
	// 	fmt.Fprintln(tw)
	// 	fmt.Fprintln(tw)
	// 	tw.Flush()
	// }
	// color.Unset()

	println("Starting HTTP Server")
	go runHTTPServer()

	// Run Go routine in later
	// Running in go now will close the main thread
	RunRPCServer()

	println("Starting RPC Server")

}
