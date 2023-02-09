package database

import (
	"fmt"
	"os"

	"github.com/dilaragorum/ticket-api/internal/ticket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Setup() (*gorm.DB, error) {
	constr := fmt.Sprintf(
		"host=%s port=%s  user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	var err error
	db, err = gorm.Open(postgres.Open(constr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate() {
	db.AutoMigrate(&ticket.Ticket{})   //nolint:errcheck
	db.AutoMigrate(&ticket.Purchase{}) //nolint:errcheck
}
