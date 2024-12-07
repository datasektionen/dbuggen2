package database

import (
	"database/sql"
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
										WHERE id IS NOT NULL AND
											publishing_date IS NOT NULL
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
									WHERE id IS NOT NULL AND
										publishing_date IS NOT NULL
									ORDER BY publishing_date DESC`)

		if err != nil {
			log.Println(err)
			return issues, err
		}
	}

	return issues, nil
}

// Gets all articles in a certain issue. Will return an error if any article
// is not nØllesafe.
func GetArticles(db *sqlx.DB, issue int, darkmode bool) ([]Article, error) {
	var articles []Article

	if err := db.Select(&articles, `SELECT * FROM Archive.Article
										WHERE issue=$1
											ORDER BY issue_index`,
		issue); err != nil {
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

// Creates a list of all authors who've contributed to an issue. Lists them in
// the order of which articles they've written.
func GetAuthorsForIssue(db *sqlx.DB, issueID int) ([][]Author, error) {
	type authoredArticle struct {
		IssueIndex   int            `db:"issue_index"`
		KthID        string         `db:"kth_id"`
		PreferedName sql.NullString `db:"prefered_name"`
	}

	var authoredArticles []authoredArticle
	err := db.Select(&authoredArticles, `SELECT issue_index, kth_id, prefered_name FROM (
											Archive.Member FULL JOIN (
												Archive.Article FULL JOIN Archive.AuthoredBy ON
												Archive.Article.id = Archive.AuthoredBy.article_id)
												USING (kth_id))
											WHERE issue=$1 AND kth_id IS NOT NULL
											ORDER BY issue_index ASC`, issueID)

	if err != nil {
		log.Println(err)

		var a [][]Author
		return a, err
	}

	if len(authoredArticles) == 0 {
		var a [][]Author
		return a, nil
	}

	authors := make([][]Author, authoredArticles[len(authoredArticles)-1].IssueIndex+1)

	for _, a := range authoredArticles {
		authors[a.IssueIndex] = append(authors[a.IssueIndex], Author{a.KthID, a.PreferedName})
	}

	return authors, nil
}

func GetAuthorsForArticle(db *sqlx.DB, article int) ([]Author, error) {
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

func GetActiveMembers(db *sqlx.DB) ([]Member, error) {
	var members []Member
	err := db.Select(&members, `SELECT kth_id, prefered_name, hosted_url, title, active
									FROM (Archive.Member FULL JOIN (
											SELECT id AS picture, hosted_url
												FROM Archive.External
												WHERE type_of_external = 'image'
											) AS ext USING(picture))
										WHERE active = true`)

	if err != nil {
		log.Println(err)
		return members, err
	}

	return members, nil
}

func UpdateIssue(db *sqlx.DB, articlesModify []Article, articlesNew []Article, articlesRemove []int) error {
	tx := db.MustBegin()
	for _, a := range articlesNew {
		tx.MustExec(`INSERT INTO Archive.Article (
							title,
							issue,
							author_text,
							issue_index,
							content,
							last_edited,
							n0lle_safe
						) VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, a.Title, a.Issue, a.AuthorText, a.IssueIndex,
			a.Content, a.LastEdited, a.N0lleSafe)
	}

	for _, a := range articlesModify {
		tx.MustExec(`UPDATE Archive.Article SET
							title=$2,
							author_text=$3,
							issue_index=$4,
							content=$5,
							last_edited=$6,
							n0lle_safe=$7
						WHERE id=$1
		`, a.ID, a.Title, a.AuthorText, a.IssueIndex,
			a.Content, a.LastEdited, a.N0lleSafe)
	}

	for _, id := range articlesRemove {
		tx.MustExec(`DELETE FROM Archive.Article WHERE id=$1`, id)
	}

	return tx.Commit()
}

func DeleteIssue(db *sqlx.DB, issue int) error {
	tx := db.MustBegin()
	tx.MustExec(`DELETE FROM Archive.Issue where id=$1`, issue)
	return tx.Commit()
}
