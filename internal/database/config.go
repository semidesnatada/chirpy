package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type State struct {
	DB *Queries
}

func CreateState() State {
	
	godotenv.Load()
	DB_URL := os.Getenv("DB_URL")
	
	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		fmt.Printf("Could not connect to database: %s\n", err.Error())
		os.Exit(1)
	}
	dbQueries := New(db)

	return State{DB: dbQueries}

}