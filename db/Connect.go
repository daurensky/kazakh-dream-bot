package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func Connect() (*sql.DB, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	dbClient, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, err
	}

	return dbClient, nil
}
