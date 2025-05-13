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
	JWTSecret string
	PolkaSecret string
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

	JWT_secret := os.Getenv("JWT_SECRET")
	if JWT_secret == "" {
		fmt.Println("JWT_SECRET must be set")
		os.Exit(1)
	}

	Polka_secret := os.Getenv("POLKA_KEY")
	if Polka_secret == "" {
		fmt.Println("POLKA_KEY must be set")
		os.Exit(1)
	}

	return apiConfig{
		fileserverHits: atomic.Int32{},
		DB: dbQueries,
		JWTSecret: JWT_secret,
		PolkaSecret: Polka_secret,
	}

}