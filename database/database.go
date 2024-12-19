package database

import (
	"log"
	"os"
	"payment-gateway/pkg/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func InitializeDB() {
	uri := os.Getenv("POSTGRES_URI")
	if uri == "" {
		log.Fatalf("POSTGRES_URI environment variable is not set")
	}

	var err error

	err = utils.RetryOperation(func() error {
		db, err = sqlx.Open("postgres", uri)
		if err != nil {
			return err
		}
		db.SetConnMaxLifetime(300 * time.Second)
		db.SetMaxOpenConns(50)
		db.SetMaxIdleConns(100)

		return db.Ping()
	}, 5)

	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	log.Println("Postgres is Connected")
}

func GetDB() *sqlx.DB {
	return db
}
