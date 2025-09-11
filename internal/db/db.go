package db

import (
	"database/sql"
	"log"
)

func Connect(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("could not ping db: %v", err)
	}

	return db
}
