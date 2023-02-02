package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"text/tabwriter"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	color "github.com/fatih/color"
	factory "github.com/storage-node-p1/contracts"

	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
)

var EClient *ethclient.Client
var FS FileService
var JFS JWTFileService
var CG *ContractGateway
var NC *NodeConfig
var GRPCServer *grpc.Server
var DBS *DBService
var NodeWallet *Wallet

func initContract(contractAddress string) {
	// color.Set(color.FgYellow)
	var cgErr error
	CG, cgErr = NewContractRPC(contractAddress)
	if cgErr != nil {
		color.Set(color.FgRed)
		log.Fatalln(cgErr)

	}

	log.Println("Fetching Contract Stats...")
	color.Set(color.FgGreen)
	contractStats, err := CG.GetContractStats()
	if err != nil {
		color.Set(color.FgRed)
		// log.Println("Failed to fetch contract stats")
		log.Println(err)
		panic("Failed to establish connection with contract.")
		// color.Unset()
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

	printLn("Opening DB")
	var err error
	DBS, err = connectDB()
	if err != nil {
		panic(err)
	}
}

func printLn(v ...interface{}) {
	log.Print(v...)
}

func deployNewContract() (string, error) {
	color.Set(color.FgBlue)
	acc := NodeWallet.account
	addr := acc.Address
	blnc, err := EClient.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		printLn("Failed to fetch balance")
		panic(err)
	}
	fmt.Printf("Your wallet %v with %v eth balance will be used. Continue? (y/N): ", addr.Hex(), weiToEther(blnc))
	var res string
	fmt.Scanln(&res)
	color.Unset()
	if (res != "y" && res == "Y") || len(res) == 0 {
		return "", errors.New("operation cancelled")
	}

	gassPrice, err := EClient.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	chainId, err := EClient.ChainID(context.Background())
	if err != nil {
		printLn("Failed to fetch chain id")
		return "", err
	}
	auth, err := bind.NewKeyStoreTransactorWithChainID(NodeWallet.ks, *NodeWallet.account, chainId)
	if err != nil {
		return "", err
	}
	auth.Value = big.NewInt(0)
	// auth.GasLimit = 0 // estimate
	auth.GasPrice = gassPrice
	// nonce, err := EClient.PendingNonceAt(context.Background(), addr)
	if err != nil {
		return "", err
	}
	// auth.Nonce = big.NewInt(int64(nonce))
	auth.NoSend = true
	if NC.FactoryAddress == "" {
		return "", errors.New("factoryAddress is not specified in config file")
	}

	println("Deploying contract...")
	fact, err := factory.NewStorageNodeFactoryContract(common.HexToAddress(NC.FactoryAddress), EClient)
	if err != nil {
		printLn(err)
		return "", errors.New("failed to initialize factory contract")
	}
	tx, err := fact.CreateStorageContract(auth, []byte{}, NC.RpcURL)
	if err != nil {
		printLn(err)
		return "", errors.New("failed to create new storage contract with factory")
	}

	color.Set(color.FgBlue)
	fmt.Printf("Estimations are as follows: Gas Limit = %v, Gass Price = %v, Tx Cost = %v eth. Continue? (y/N): ", tx.Gas(), tx.GasPrice(), weiToEther(tx.Cost()))
	fmt.Scanln(&res)

	err = EClient.SendTransaction(context.Background(), tx)

	if err != nil {
		color.Set(color.FgRed)
		printLn("Failed to send tx")
		printLn(err)
		color.Unset()
		return "", nil
	}

	printLn("Waiting for transaction to be mined...")
	rct, err := bind.WaitMined(context.Background(), EClient, tx)
	if err != nil {
		printLn(err)
		return "", errors.New("transaction mining failed")
	}
	factoryABI, err := abi.JSON(strings.NewReader(factory.StorageNodeFactoryContractABI))
	if err != nil {
		printLn(err)
		return "", errors.New("failed to parse factory ABI")
	}
	var event factory.StorageNodeFactoryContractEvNewStorageContract

	if len(rct.Logs) == 0 {
		return "", errors.New("tx didn't emit any contract deployment event")
	}

	err = factoryABI.UnpackIntoInterface(&event, "EvNewStorageContract", rct.Logs[0].Data)
	if err != nil {
		printLn(err)
		return "", errors.New("failed to unpack event data")
	}
	if len(event.Addr.String()) == 0 {
		return "", errors.New("got empty address in event")
	}
	address := event.Addr

	color.Set(color.FgGreen)
	printLn("Deployed contract address: " + addr.Hex())
	color.Unset()

	return address.Hex(), nil

}
func main() {
	// mainMerkle()

	// proof, dir := GenerateMerkleProof(1, [][]byte{
	// 	[]byte("hello"),
	// 	[]byte("world!"),
	// 	[]byte("hello"),
	// 	[]byte("hello"),
	// 	[]byte("hello"),
	// })
	// printLn(proof)
	// printLn(dir)

	// panic("exit")
	log.SetFlags(log.Ldate | log.Ltime)
	printLn("Booting...")
	openDB()
	NodeWallet = LoadWallet()
	color.Set(color.FgYellow)
	printLn("Using " + NodeWallet.account.Address.Hex() + " as signer.")
	color.Unset()

	printLn("Loading Configuration")

	NC = LoadConfig()

	gateway := NC.ProviderURL // Todo: Replace with CONFIG VAR
	ec, err := ethclient.Dial(gateway)

	if err != nil {
		color.Set(color.FgRed)
		printLn("Failed to connect to " + gateway + " gateway")
		panic(err)
	} else {
		EClient = ec
		color.Set(color.FgYellow)
		printLn("Attached client to " + gateway + " gateway")
		color.Unset()
	}
	contractAddr := ""
	if len(NC.ContractAddress) == 0 {
		color.Set(color.FgBlue, color.Bold)
		// color.Set(color.FgGreen,)
		fmt.Printf("Contract address is not present in configuration. Would you like to deploy new contract? (y/N): ")
		var c string
		fmt.Scanln(&c)
		color.Unset()

		if c == "" || (c != "y" && c != "Y") {
			color.Set(color.FgYellow)
			printLn("Active contract is required to proceed further. Exiting")
			os.Exit(0)
		} else {
			contractAddr, err = deployNewContract()
			if err != nil {
				color.Set(color.FgYellow)
				printLn(err)
				printLn("Active contract is required to proceed further. Exiting")
				os.Exit(0)
			}

			printLn("Saving contract address in config file.")
			NC.ContractAddress = contractAddr

			configBytes, err := yaml.Marshal(NC)
			if err != nil {
				panic(err)
			}
			err = os.WriteFile("config.yaml", configBytes, 0666)
			if err != nil {
				panic(err)
			}
		}
	} else {
		contractAddr = NC.ContractAddress
		color.Set(color.FgYellow)

		printLn("Using contract from configuration " + contractAddr)
		color.Unset()
	}

	initContract(contractAddr)

	color.Set(color.FgGreen)
	println("Total Storage:", NC.TotalStorage, "B")
	color.Unset()

	chainId, err := EClient.ChainID(context.Background())
	if err != nil {
		panic("Failed to fetch chain id")
	}
	tx, err := bind.NewKeyStoreTransactorWithChainID(NodeWallet.ks, *NodeWallet.account, chainId)

	if err != nil {
		panic("Failed to get signer")
	}

	blnc, err := EClient.BalanceAt(context.Background(), CG.contractAddress, nil)
	if err == nil {
		printLn("Balance: ", blnc)

		collateral, err := CG.instance.LockedCollateral(nil)
		if err != nil {
			println("LockedCollateral error")
			printLn(err)
		} else {
			printLn("Collateral: ", collateral)
			remaining := blnc.Sub(blnc, collateral)
			printLn("Remaining: ", remaining)

			if remaining.Cmp(big.NewInt(0)) == 1 {
				_, err = CG.instance.Withdraw(tx, remaining, NodeWallet.account.Address)
				if err != nil {
					printLn("Withdraw error")
					printLn(err)
				} else {
					printLn("Withdraw: ", remaining)
				}
			}
		}
	}

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

	InitSchedular()
	printLn("Starting HTTP Server")
	go runHTTPServer()

	// Run Go routine in later
	// Running in go now will close the main thread
	printLn("Starting RPC Server")
	RunRPCServer()

}
