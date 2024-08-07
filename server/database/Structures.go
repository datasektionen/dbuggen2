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

// Relevant information for issue on home page
type HomeIssue struct {
	ID             int
	Title          string
	PublishingDate time.Time `db:"publishing_date"`
	Coverpage      sql.NullString
	Views          int
}

type Article struct {
	ID         int
	Title      string
	Issue      int
	AuthorText sql.NullString `db:"author_text"`
	IssueIndex int            `db:"issue_index"`
	Content    string
	LastEdited time.Time `db:"last_edited"`
	N0lleSafe  bool      `db:"n0lle_safe"`
}

type Author struct {
	KthID        string         `db:"kth_id"`
	PreferedName sql.NullString `db:"prefered_name"`
}
