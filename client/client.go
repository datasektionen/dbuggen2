package client

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"dbuggen/server/database"
)

// Home page
func Home(db *sqlx.DB, ds *DarkmodeStatus) func(c *gin.Context) {
	return func(c *gin.Context) {
		issuesRaw, err := database.GetHomeIssues(db, Darkmode(ds))
		if err != nil {
			c.Redirect(http.StatusInternalServerError, "")
			return
		}

		type DisplayIssue struct {
			IssueID        string
			Title          string
			PublishingDate string
			Coverpage      template.HTML
			Views          int
		}

		var issues []DisplayIssue
		for _, iss := range issuesRaw {
			issues = append(issues,
				DisplayIssue{
					fmt.Sprintf("issue/%v", iss.ID),
					iss.Title,
					iss.PublishingDate.Format(time.DateOnly),
					coverpage(iss.Coverpage),
					iss.Views})
		}
		c.HTML(http.StatusOK, "home.html", gin.H{
			"pagetitle": "dbuggen",
			"issues":    issues,
		})
	}
}

// Abritrary issue featuring all the articles
func Issue(db *sqlx.DB, ds *DarkmodeStatus) func(c *gin.Context) {
	type issueArticle struct {
		Title       string
		ArticleLink string
		Authors     string
		Content     template.HTML
		LastEdited  string
	}

	return func(c *gin.Context) {
		issueID, err := pathIntSeparator(c.Param("issue"))
		if err != nil {
			c.Redirect(http.StatusBadRequest, "")
			return
		}

		darkmode := Darkmode(ds)

		issue, err := database.GetIssue(db, issueID, darkmode)
		if err != nil {
			c.Redirect(http.StatusInternalServerError, "")
			return
		}

		articles, err := database.GetArticles(db, issueID, darkmode)
		if err != nil {
			c.Redirect(http.StatusInternalServerError, "")
			return
		}

		var issueArticles []issueArticle
		for _, article := range articles {
			var authors string
			if article.AuthorText.Valid {
				emptyAuthors := make([]database.Author, 0)
				authors = authortext(article.AuthorText, emptyAuthors)
			} else {
				databaseAuthors, err := database.GetAuthors(db, article.ID)
				if err != nil {
					c.Redirect(http.StatusInternalServerError, "")
					return
				}

				authors = authortext(article.AuthorText, databaseAuthors)
			}

			content := mdToHTML(article.Content)
			lastEdited := article.LastEdited.Format(time.DateOnly)
			issueArticle := issueArticle{
				Title:       article.Title,
				ArticleLink: fmt.Sprintf("/issue/%v/%v", issue.ID, article.IssueIndex),
				Authors:     authors,
				Content:     content,
				LastEdited:  lastEdited,
			}
			issueArticles = append(issueArticles, issueArticle)
		}

		c.HTML(http.StatusOK, "issue.html", gin.H{
			"coverpage":  coverpage(issue.Coverpage),
			"issueTitle": issue.Title,
			"articles":   issueArticles,
		})
	}
}

// Arbitrary article
func Article(db *sqlx.DB, ds *DarkmodeStatus) func(c *gin.Context) {
	return func(c *gin.Context) {
		issueID, errI := pathIntSeparator(c.Param("issue"))
		articleIndex, errA := pathIntSeparator(c.Param("article"))
		if errI != nil || errA != nil {
			c.Redirect(http.StatusBadRequest, "")
			return
		}

		article, err := database.GetArticle(db, issueID, articleIndex, Darkmode(ds))
		if err != nil {
			c.Redirect(http.StatusInternalServerError, "")
			return
		}
		authors, err := database.GetAuthors(db, article.ID)
		if err != nil {
			c.Redirect(http.StatusInternalServerError, "")
			return
		}

		c.HTML(http.StatusOK, "article.html", gin.H{
			"pagetitle":      article.Title,
			"title":          article.Title,
			"authors":        authortext(article.AuthorText, authors),
			"articleContent": mdToHTML(article.Content),
		})
	}
}
