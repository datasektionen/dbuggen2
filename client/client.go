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
					fmt.Sprintf("issue/%v/0", iss.ID),
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

// Page for all of (active) redaqtionen to be shown to the world
func Redaqtionen(db *sqlx.DB, DFUNKT_URL string) func(c *gin.Context) {
	return func(c *gin.Context) {
		members, err := database.GetActiveMembers(db)
		if err != nil {
			c.Redirect(http.StatusInternalServerError, "")
			return
		}

		chefredIDs := getChefreds(DFUNKT_URL)
		chefreds, members := removeDuplicateChefreds(chefredIDs, members)
		displaymembers := displaymemberize(members)
		displayChefreds := displaymemberize(chefreds)

		c.HTML(http.StatusOK, "redaqtionen.html", gin.H{
			"chefreds": displayChefreds,
			"members":  displaymembers,
		})
	}
}
