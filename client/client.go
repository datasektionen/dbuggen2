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
func Home(issuesRaw []database.HomeIssue) func(c *gin.Context) {
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
			DisplayIssue{fmt.Sprintf("issue/%v/0", iss.ID),
				iss.Title,
				iss.PublishingDate.Format(time.DateOnly),
				coverpage(iss.Coverpage),
				iss.Views})
	}

	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.tmpl", gin.H{
			"pagetitle": "dbuggen",
			"issues":    issues,
		})
	}
}

// Arbitrary article
func Article(db *sqlx.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		issueID, errI := pathIntSeparator(c.Param("issue"))
		articleIndex, errA := pathIntSeparator(c.Param("article"))
		if errI != nil || errA != nil {
			c.Redirect(http.StatusBadRequest, "")
		}

		article := database.GetArticle(db, issueID, articleIndex)
		authors := database.GetAuthors(db, article.ID)

		c.HTML(http.StatusOK, "article.tmpl", gin.H{
			"pagetitle":      article.Title,
			"title":          article.Title,
			"authors":        authortext(article.AuthorText, authors),
			"articleContent": mdToHTML(article.Content),
		})
	}
}
