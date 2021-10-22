package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Sql struct {
	Db *sqlx.DB
}

var Sqlx *Sql

// Setup initializes the database instance
func Setup() {
	Sqlx = &Sql{}

	dataSource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "admin", "scrap")
	Sqlx.Db = sqlx.MustConnect("postgres", dataSource)
	if err := Sqlx.Db.Ping(); err != nil {
		log.Fatalf("db.Setup err: %v", err)
		Sqlx.Db.Close()
	}
}
