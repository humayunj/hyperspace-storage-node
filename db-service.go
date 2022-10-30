package main

import (
	"database/sql"
	"log"

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

func seedDB() {
	const sql = `
		CREATE TABLE transactions(
			id VARCHAR(32)
			user_address VARCHAR(20)
			file_merkle_hash VARCHAR(32)
			file_name VARCHAR(50),
			file_size INT UNSIGNED,
			
			uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			
			status INT, -- CANCELLED, COMPLETE, IN_PROGRESS
			expires_at DATETIME,
			bid_price string,

		);
	`
}

func (dbs *DBService) connect() (*DBService, error) {
	db, err := sql.Open("sqlite3", "data.db")

	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}
	log.Println("DB Connected")

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
