package belajar_golang_database

import (
	"database/sql"
	"time"
)

func GetConnection() *sql.DB {

	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/belajar_golang_database?parseTime=true")
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)                  // set minimum connection created on first
	db.SetMaxOpenConns(100)                 // set maximum connection can created
	db.SetConnMaxIdleTime(5 * time.Minute)  // set max time connection not working anything
	db.SetConnMaxLifetime(60 * time.Minute) // set max time connection must be renew

	return db

}
