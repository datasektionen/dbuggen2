package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Start(db_url string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func GetIssues(db *sqlx.DB) []Issue {
	issues := []Issue{}

	err := db.Select(&issues, "SELECT * FROM Archive.Issue")
	if err != nil {
		log.Fatal(err)
	}

	return issues
}
