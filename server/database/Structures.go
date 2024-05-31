package database

import (
	"database/sql"
	"time"
)

type Member struct {
	Kth_id string
	Title  string
	Active bool
}

type Issue struct {
	ID             int
	Title          string
	PublishingDate time.Time `db:"publishing_date"`
	Pdf            sql.NullInt32
	Html           sql.NullInt32
	Coverpage      sql.NullInt32
	Views          int
}
