package database

import (
	"database/sql"
	"log"
	"time"
)

// DbConn : the reference for the db connection
var DbConn *sql.DB

// SetupDatabase : sets up the db connection
func SetupDatabase() {
	var err error
	DbConn, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/inventorydb")
	if err != nil {
		log.Fatal(err)
	}
	DbConn.SetMaxOpenConns(4)
	DbConn.SetMaxIdleConns(4)
	DbConn.SetConnMaxLifetime(60 * time.Second)
}
