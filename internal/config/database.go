package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Fallback/Default for development if not set, slightly dangerous but helpful for local test if they forget .env
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if dbname == "" {
		dbname = "pdf_management"
	}
	// Password usually must be set, but let's see.

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Printf("Warning: Could not connect to database: %v. Make sure DB is running and credentials are correct.", err)
		// We don't Fatal here to allow main to run and maybe print help, but for a real backend we might want to fail.
		// Actually, let's Fatal because per requirements "Connection to DB" is critical.
		log.Fatal(err)
	}

	log.Println("Successfully connected to the database")
}
