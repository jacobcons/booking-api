package db

import (
	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sqlx.DB

func init() {
	var err error
	DB, err = sqlx.Connect("postgres", os.Getenv("DBSTRING"))
	DB.MapperFunc(strcase.ToSnake)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
}
