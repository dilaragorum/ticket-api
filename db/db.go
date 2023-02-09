package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" //nolint:revive,nolintlint
)

func Connect() (*sql.DB, error) {
	connectionURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_DB"))

	fmt.Println(connectionURL)

	db, err := sql.Open("postgres", connectionURL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
