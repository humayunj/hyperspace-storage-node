package main

import (
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type Wallet struct {
	ks      *keystore.KeyStore
	account *accounts.Account
}

func LoadWallet() *Wallet {
	// file := "./store/keystore"
	ks := keystore.NewKeyStore("./store", keystore.StandardScryptN, keystore.StandardScryptP)

	password := os.Getenv("store_password")
	if len(password) == 0 {
		panic("store_password env var is not set")
	}

	if len(ks.Accounts()) == 0 {
		panic("keystore file not present")
	}
	account := ks.Accounts()[0]
	err := ks.Unlock(account, password)
	if err != nil {
		panic(err)
	}
	return &Wallet{
		ks:      ks,
		account: &account,
	}
}
