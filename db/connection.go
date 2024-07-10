package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sqlx.DB

func init() {
	var err error
	DB, err = sqlx.Connect("postgres", os.Getenv("DBSTRING"))
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
}
