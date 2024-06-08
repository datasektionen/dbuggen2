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

func GetIssues(db *sqlx.DB) ([]Issue, error) {
	issues := []Issue{}

	err := db.Select(&issues, "SELECT * FROM Archive.Issue ORDER BY publishing_date DESC")
	if err != nil {
		log.Println(err)
		return issues, err
	}

	return issues, nil
}

// haha.
func GetHomeIssues(db *sqlx.DB, darkmode bool) ([]HomeIssue, error) {
	issues := []HomeIssue{}

	if darkmode { // if the mörkläggning is active
		err := db.Select(&issues, `WITH safe_issues AS (
										SELECT * FROM Archive.Issue
											WHERE id IN
												(SELECT issue FROM Archive.Article
													WHERE n0lle_safe = TRUE)
									)
									SELECT id, title, publishing_date, hosted_url AS coverpage, views
										FROM (safe_issues FULL JOIN (
											SELECT id AS coverpage, hosted_url
												FROM Archive.External
												WHERE type_of_external = 'image'
											) AS ext
											USING(coverpage))
										WHERE id IS NOT NULL
										ORDER BY publishing_date DESC`)

		if err != nil {
			log.Println(err)
			return issues, err
		}
	} else {
		err := db.Select(&issues, `SELECT id, title, publishing_date, hosted_url AS coverpage, views
									FROM (Archive.Issue FULL JOIN (
										SELECT id AS coverpage, hosted_url
											FROM Archive.External
											WHERE type_of_external = 'image'
										) AS ext
										USING(coverpage))
									WHERE id IS NOT NULL
									ORDER BY publishing_date DESC`)

		if err != nil {
			log.Println(err)
			return issues, err
		}
	}

	return issues, nil
}

// GetArticles retrieves a list of article IDs from the database for a given issue.
func GetArticles(db *sqlx.DB, issue int) ([]int, error) {
	var articles []int
	err := db.Get(&articles, `SELECT id FROM Archive.Article WHERE issue=$1 ORDER BY issue_index ASC`, issue)
	if err != nil {
		log.Println(err)
		return articles, err
	}

	return articles, nil
}

func GetArticle(db *sqlx.DB, issueID int, index int, darkmode bool) (Article, error) {
	var article Article

	if darkmode {
		if err := db.Get(&article, `SELECT * FROM Archive.Article
										WHERE issue=$1
											AND issue_index=$2
											AND issue IN (
												SELECT id FROM Archive.Issue
													WHERE id IN (
														SELECT issue FROM Archive.Article
															WHERE n0lle_safe = TRUE))`, issueID, index); err != nil {
			log.Println(err)
			return article, err
		}
	} else {
		if err := db.Get(&article, "SELECT * FROM Archive.Article WHERE issue=$1 AND issue_index=$2", issueID, index); err != nil {
			log.Println(err)
			return article, err
		}
	}

	return article, nil
}

func GetAuthors(db *sqlx.DB, article int) ([]Author, error) {
	var authors []Author
	err := db.Select(&authors, `SELECT kth_id, prefered_name FROM
								(Archive.Member LEFT JOIN Archive.AuthoredBy USING(kth_id))
								WHERE article_id=$1`, article)
	if err != nil {
		log.Println(err)
		return authors, err
	}

	return authors, nil
}
