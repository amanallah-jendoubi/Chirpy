package config

import (
	"database/sql"
	"fmt"
	"github.com/amanallah-jendoubi/Textio/http/helpers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sql.DB //pool manager

func Init(envPath string) {
	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("Warning: could not load .env file from %s: %v", envPath, err)
		}
	}
	host := helpers.GetEnv("POSTGRES_HOST", "db")
	port := helpers.GetEnv("POSTGRES_PORT", "5432")
	user := helpers.GetEnv("POSTGRES_USER", "postgres")
	dbname := helpers.GetEnv("POSTGRES_DB", "postgres")
	password := os.Getenv("POSTGRES_PASSWORD")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname) //no ssl fo db connections (dev)

	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Warning: failed to close existing DB connection: %v", err)
		}
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}

	if err = DB.Ping(); err != nil {
		DB.Close()
		log.Fatalf("Failed to ping DB: %v", err)
	}
}
