package database

import (
	"errors"
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

func GetIssue(db *sqlx.DB, issueID int, darkmode bool) (HomeIssue, error) {
	var issue HomeIssue

	if darkmode {
		err := db.Get(&issue, `WITH safe_issues AS (
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
									WHERE id=$1`, issueID)
		if err != nil {
			log.Println(err)
			return issue, err
		}
	} else {
		err := db.Get(&issue, `SELECT id, title, publishing_date, hosted_url AS coverpage, views
									FROM (Archive.Issue FULL JOIN (
										SELECT id AS coverpage, hosted_url
											FROM Archive.External
											WHERE type_of_external = 'image'
										) AS ext
										USING(coverpage))
									WHERE id=$1`, issueID)

		if err != nil {
			log.Println(err)
			return issue, err
		}
	}

	return issue, nil
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
func GetArticles(db *sqlx.DB, issue int, darkmode bool) ([]Article, error) {
	var articles []Article

	if err := db.Select(&articles, `SELECT * FROM Archive.Article WHERE issue=$1`, issue); err != nil {
		log.Println(err)
		return articles, err
	}

	if darkmode {
		for _, article := range articles {
			if !article.N0lleSafe {
				return articles, errors.New("not safe")
			}
		}
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
