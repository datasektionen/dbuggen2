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

	err := db.Select(&issues, "SELECT * FROM Archive.Issue ORDER BY publishing_date DESC")
	if err != nil {
		log.Fatal(err)
	}

	return issues
}

// GetArticles retrieves a list of article IDs from the database for a given issue.
func GetArticles(db *sqlx.DB, issue int) []int {
	var articles []int
	err := db.Get(&articles, `SELECT id FROM Archive.Article WHERE issue=$1 ORDER BY issue_index ASC`, issue)
	if err != nil {
		log.Fatal(err)
	}

	return articles
}

func GetArticleFromID(db *sqlx.DB, id int) Article {
	var article Article
	err := db.Get(&article, "SELECT * FROM Archive.Article WHERE id=$1", id)
	if err != nil {
		log.Fatal(err)
	}

	return article
}

func GetArticle(db *sqlx.DB, issueID int, index int) Article {
	var article Article
	err := db.Get(&article, "SELECT * FROM Archive.Article WHERE issue=$1 AND issue_index=$2", issueID, index)
	if err != nil {
		log.Fatal(err)
	}

	return article
}
