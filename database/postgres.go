package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

// ConnectDB menginisialisasi koneksi ke database
func ConnectDB() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable not set")
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Could not ping the database: %v", err)
	}

	log.Println("Successfully connected to the database")
}