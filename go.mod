module github.com/storage-node-p1

go 1.16

replace github.com/storage-node-p1/storage-node-proto => ./proto/

replace github.com/storage-node-p1 => ./

replace github.com/storage-node-p1/contracts => ./contracts/

require (
	github.com/btcsuite/btcd/btcec/v2 v2.3.0 // indirect
	github.com/cbergoon/merkletree v0.2.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.1.0 // indirect
	github.com/dgraph-io/badger/v3 v3.2103.4
	github.com/ethereum/go-ethereum v1.10.26
	github.com/fatih/color v1.13.0
	github.com/go-co-op/gocron v1.17.1
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/storage-node-p1/contracts v0.0.0-00010101000000-000000000000
	github.com/storage-node-p1/storage-node-proto v0.0.0-00010101000000-000000000000
	github.com/wealdtech/go-merkletree v1.0.0
	golang.org/x/crypto v0.1.0 // indirect
	google.golang.org/grpc v1.50.1
	gopkg.in/yaml.v3 v3.0.1

)
