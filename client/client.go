package client

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	db "dbuggen/server/database"
)

// Home page
func Home(issuesRaw []db.Issue) func(c *gin.Context) {
	type RelevantIssue struct {
		Name       string
		Date       time.Time
		Coverimage string
	}

	var issues []RelevantIssue
	for _, iss := range issuesRaw {
		issues = append(issues, RelevantIssue{iss.Title, iss.PublishingDate, "https://dbuggen.s3.eu-west-1.amazonaws.com/dbuggen2/marke.png"})
	}

	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":   "Main wow",
			"content": "home.tmpl",
			// Specific for home.tmpl ↓
			"issues": issues,
		})
	}
}
