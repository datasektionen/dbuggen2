package client

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"dbuggen/server/database"
	db "dbuggen/server/database"
)

// Home page
func Home(issuesRaw []db.Issue) func(c *gin.Context) {
	type RelevantIssue struct {
		Title          string
		PublishingDate time.Time
		Coverimage     string
	}

	var issues []RelevantIssue
	for _, iss := range issuesRaw {
		issues = append(issues, RelevantIssue{iss.Title, iss.PublishingDate, "https://dbuggen.s3.eu-west-1.amazonaws.com/dbuggen2/marke.png"})
	}

	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"pagetitle": "dbuggen",
			"content":   "hcontent",
			// Specific for home.tmpl ↓
			"issues": issues,
		})
	}
}

// Arbitrary article
func Article(db *sqlx.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		issueID := pathIntSeparator(c.Param("issue"))
		articleIndex := pathIntSeparator(c.Param("article"))

		article := database.GetArticle(db, issueID, articleIndex)
		fmt.Println(article.Content)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"pagetitle": article.Title,
			"content":   "acontent",
			// Specific for article.tmpl ↓
			"title":          article.Title,
			"authors":        article.AuthorText,
			"articleContent": article.Content,
		})
	}
}

// Function to remove the "/" before parameters, which was
// a problem. Turns "/123": string, into 123: int.
func pathIntSeparator(paramRaw string) int {
	_, paramLessRaw, found := strings.Cut(paramRaw, "/")

	if found {
		param, err := strconv.Atoi(paramLessRaw)
		if err != nil {
			log.Fatal(err)
		}
		return param
	}

	param, err := strconv.Atoi(paramRaw)
	if err != nil {
		log.Fatal(err)
	}
	return param
}
