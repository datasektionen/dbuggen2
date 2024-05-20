package database

import (
	"database/sql"
	"fmt"
	"log"

	// "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Member struct {
	kth_id string
	title  string
	active bool
}

func DatabaseDo(db_url string) {
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if db.Ping() != nil {
		fmt.Println("AAAAAAAAHH!")
	}
}
