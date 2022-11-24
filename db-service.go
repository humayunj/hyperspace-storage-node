package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type DBService struct {
	db *sql.DB
}

// Possible routines
/**

- getUserFiles(address)
- getFileDetails(fileId)
- addFile()
- removeFile()

*/

type TransactionStatus uint8

const (
	TRANSACTION_STATUS_IN_PROGRESS TransactionStatus = 0
	TRANSACTION_STATUS_COMPLETE    TransactionStatus = 1
	TRANSACTION_STATUS_CANCELLED   TransactionStatus = 2
	TRANSACTION_STATUS_PENDING     TransactionStatus = 3
)

type TransactionParams struct {
	FileKey            string
	UserAddress        string
	FileMerkleRootHash string
	FileName           string
	FileSize           uint64
	Status             TransactionStatus
	BidPrice           string
	ExpiresAt          uint64
	UploadedAt         uint64
}

func (dbs *DBService) InsertTransaction(params TransactionParams) error {

	_, err := dbs.db.Exec(
		`INSERT INTO transactions 
			(
				file_key,user_address,file_merkle_hash,
				file_name,file_size,status,expires_at,
				bid_price,uploaded_at
			) 	
				VALUES(?,?,?,?,?,?,?,?,?)`,
		params.FileKey,
		params.UserAddress,
		params.FileMerkleRootHash,
		params.FileName,
		params.FileSize,
		params.Status,
		params.ExpiresAt,
		params.BidPrice,
		params.UploadedAt,
	)
	return err
}

func (dbs *DBService) GetTransaction(fileKey string) (TransactionParams, error) {

	var tp TransactionParams

	err := dbs.db.QueryRow(
		`SELECT file_key,user_address,file_merkle_hash,
		file_name,file_size,status,expires_at,
		bid_price,uploaded_at FROM  transactions  WHERE file_key = ?`,
		fileKey,
	).Scan(&tp.FileKey, &tp.UserAddress,
		&tp.FileMerkleRootHash, &tp.FileName,
		&tp.FileSize, &tp.Status, &tp.ExpiresAt,
		&tp.BidPrice, &tp.UploadedAt)
	return tp, err
}
func seedDB(db *sql.DB) error {
	const query = `
		CREATE TABLE IF NOT EXISTS transactions (
			file_key VARCHAR(256) PRIMARY KEY UNIQUE,
			user_address VARCHAR(256),
			file_merkle_hash VARCHAR(256),
			file_name VARCHAR(50),
			file_size INT UNSIGNED,
			
			status INT, -- CANCELLED, COMPLETE, IN_PROGRESS
			expires_at INT UNSIGNED,
			bid_price string,

			uploaded_at INT UNSIGNED
		);
	`

	_, err := db.Exec(query)
	if err != nil {
		color.Set(color.FgRed)
		printLn(err.Error())
		color.Unset()
		return err
	}
	printLn("Seeded DB")
	return nil
}

func connectDB() (*DBService, error) {

	seedFlag := false
	if _, err := os.Stat("data.db"); errors.Is(err, os.ErrNotExist) {
		seedFlag = true
	}

	db, err := sql.Open("sqlite3", "data.db")

	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}
	log.Println("DB Connected")

	if seedFlag {
		err = seedDB(db)
		if err != nil {
			return nil, err
		}
	}
	var d DBService
	d.db = db
	return &d, nil
	// rows, err := db.Query("SELECT * from user")
	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var name string
	// 	err = rows.Scan(&name)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	log.Println("name:", name)
	// }
}
