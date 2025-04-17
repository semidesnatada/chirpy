package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/semidesnatada/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	DB *database.Queries
}

func createState() apiConfig {
	
	godotenv.Load()
	DB_URL := os.Getenv("DB_URL")
	if DB_URL == "" {
		fmt.Println("DB_URL must be set")
		os.Exit(1)
	}
	
	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		fmt.Printf("Could not connect to database: %s\n", err.Error())
		os.Exit(1)
	}
	
	dbQueries := database.New(db)

	return apiConfig{
		fileserverHits: atomic.Int32{},
		DB: dbQueries,
	}

}