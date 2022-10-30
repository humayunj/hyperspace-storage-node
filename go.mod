module github.com/storage-node-p1

go 1.16

replace github.com/storage-node-p1/storage-node-proto => ./proto/

replace github.com/storage-node-p1 => ./

require (
	github.com/dgraph-io/badger/v3 v3.2103.2
	github.com/ethereum/go-ethereum v1.10.19
	github.com/fatih/color v1.13.0
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.1
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/storage-node-p1/storage-node-proto v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c // indirect
	google.golang.org/grpc v1.47.0
	gopkg.in/yaml.v3 v3.0.1
)
