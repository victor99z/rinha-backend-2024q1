package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Storage() *sqlx.DB {
	db := CreateConnection()
	return db
}

func CreateConnection() *sqlx.DB {

	db, err := sqlx.Connect("postgres",
		"user=admin dbname=rinha password=admin sslmode=disable host=postgres port=5432")

	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		log.Fatal(err)
		return nil
	} else {
		log.Println("Successfully Connected")
	}

	return db
}
