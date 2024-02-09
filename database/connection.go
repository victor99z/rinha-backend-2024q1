package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Storage() *sqlx.DB {
	db, err := CreateConnection()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func CreateConnection() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=admin dbname=rinha password=admin sslmode=disable host=localhost")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		log.Println("Successfully Connected")
	}

	return db, nil
}
