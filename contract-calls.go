package main

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	addrHex string,
	fileHash string,
	fileSize uint32,
	timeStart uint64,
	timeEnd uint64,
	proveTimeoutLength uint64,
	concludeTimeoutLength uint64,
	segmentsCount uint64,
	bidAmount *big.Int) error {

	addr := common.HexToAddress(addrHex)
	var hash [32]byte
	copy(hash[:], common.Hex2Bytes(fileHash)[:32])
	tx, err := cg.instance.ConcludeTransaction(nil, 0, addr, hash, fileSize, big.NewInt(time.Now().Unix()), big.NewInt(time.Now().Add(time.Minute*2).Unix()), proveTimeoutLength, concludeTimeoutLength, uint32(segmentsCount), bidAmount)
	_ = tx
	return err
}
