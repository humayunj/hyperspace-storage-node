package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func connect() {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		panic("failed to open dZ`b")
	}

	pingErr := db.Ping()
	if pingErr != nil {
		panic("Ping err")
	}
	log.Println("DB Connected")

	rows, err := db.Query("SELECT * from user")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			panic(err)
		}
		log.Println("name:", name)
	}
}
