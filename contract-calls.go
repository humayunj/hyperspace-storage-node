package main

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	contract "github.com/storage-node-p1/contracts"
)

type ContractGateway struct {
	contractAddress common.Address
	instance        *contract.StorageNodeContract
}

func NewContractRPC(contractAddress string) (*ContractGateway, error) {

	contractAddr := common.HexToAddress(contractAddress)

	instance, err := contract.NewStorageNodeContract(contractAddr, EClient)
	if err != nil {
		return nil, err
	}
	c := ContractGateway{

		contractAddr,
		instance,
	}

	return &c, nil
}

type ContractStats struct {
	LockedCollateral *big.Int
	Host             string
	TLSCert          []byte
	MappingLength    *big.Int
	Balance          *big.Int
}

func (cg *ContractGateway) GetContractStats() (*ContractStats, error) {

	host, err := cg.GetHost()
	if err != nil {
		return nil, err
	}
	cert, err := cg.GetTLSCert()
	if err != nil {
		return nil, err
	}
	coll, err := cg.GetLockedCollateral()
	if err != nil {
		return nil, err
	}

	mLen, err := cg.GetMappingsLength()
	if err != nil {
		return nil, err
	}
	blnc, err := cg.GetBalance()
	if err != nil {
		return nil, err
	}
	return &ContractStats{
		LockedCollateral: coll,
		Host:             host,
		TLSCert:          cert,
		MappingLength:    mLen,
		Balance:          blnc,
	}, nil
}

func (cg *ContractGateway) GetBalance() (*big.Int, error) {
	blnc, err := EClient.BalanceAt(context.Background(), cg.contractAddress, nil)
	if err != nil {
		return nil, err
	}
	return blnc, nil
}
func (cg *ContractGateway) GetHost() (string, error) {
	host, err := cg.instance.HOST(nil)
	return host, err
}
func (cg *ContractGateway) GetTLSCert() ([]byte, error) {
	tlsCert, err := cg.instance.TLSCert(nil)
	return (tlsCert), err
}
func (cg *ContractGateway) GetMappingsLength() (*big.Int, error) {
	mLen, err := cg.instance.MappingLength(nil)
	return mLen, err
}

func (cg *ContractGateway) GetLockedCollateral() (*big.Int, error) {
	col, err := cg.instance.LockedCollateral(nil)
	return col, err
}

func (cg *ContractGateway) ConcludeTransaction(
	userAddrHex string,
	fileHash string,
	fileSize uint32,
	timeStart uint64,
	timeEnd uint64,
	proveTimeoutLength uint64,
	concludeTimeoutLength uint64,
	segmentsCount uint64,
	bidAmount *big.Int) error {
	printLn("CLG Instance ", cg.instance == nil)

	printLn(userAddrHex, " ", fileHash, " ", fileSize, " ",
		timeStart, " ", timeEnd, " ", proveTimeoutLength, " ", concludeTimeoutLength, " ", segmentsCount, " ", bidAmount.String())
	addr := common.HexToAddress(userAddrHex)
	var hash [32]byte
	copy(hash[:], common.Hex2Bytes(fileHash)[:32])
	printLn("Addr:", addr)
	printLn("hash:", hash)
	printLn("bid:", bidAmount.String())

	host, err := cg.instance.HOST(nil)
	if err != nil {
		panic(err)
	}
	printLn("host:", host)

	if err != nil {
		return err
	}

	nonce, err := EClient.PendingNonceAt(context.Background(), NodeWallet.account.Address)
	if err != nil {
		panic(err)
	}

	chainId, err := EClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	gassPrice, err := EClient.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}
	collateralAmount := big.NewInt(1)
	collateralAmount.Mul(bidAmount, big.NewInt(2))

	// collateralAmount = collateralAmount.Add(collateralAmount, big.NewInt(10000))

	valTx := types.NewTransaction(nonce, cg.contractAddress, collateralAmount, uint64(22000), gassPrice, nil)
	sTx, err := NodeWallet.ks.SignTx(*NodeWallet.account, valTx, chainId)
	if err != nil {
		panic(err)
	}
	printLn("Sending ", collateralAmount.String(), " to contract for collateral")
	err = EClient.SendTransaction(context.Background(), sTx)

	if err != nil {
		panic(err)
	}
	_, err = bind.WaitMined(context.Background(), EClient, sTx)
	if err != nil {
		panic("error")
	}

	printLn("Mined")

	auth, err := bind.NewKeyStoreTransactorWithChainID(NodeWallet.ks, *NodeWallet.account, chainId)
	if err != nil {
		return err
	}

	auth.GasLimit = 0

	timeStartBI := new(big.Int)
	timeStartBI.SetUint64(timeStart)

	timeEndBI := new(big.Int)
	timeEndBI.SetUint64(timeEnd)

	printLn(">Conclude timeout: ", concludeTimeoutLength)
	tx, err := cg.instance.ConcludeTransaction(auth, 0,
		addr, hash, fileSize, timeStartBI, timeEndBI,
		proveTimeoutLength, concludeTimeoutLength,
		uint32(segmentsCount), bidAmount)
	_ = tx
	printLn("done")

	return err
}
