package client

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"dbuggen/server/database"
	"io"
)

// Home page
func Home(issuesRaw []database.Issue) func(c *gin.Context) {
	type RelevantIssue struct {
		IssueID        string
		Title          string
		PublishingDate time.Time
		Coverimage     string
	}

	var issues []RelevantIssue
	for _, iss := range issuesRaw {
		issues = append(issues, RelevantIssue{fmt.Sprintf("issue/%v/0", iss.ID), iss.Title, iss.PublishingDate, "https://dbuggen.s3.eu-west-1.amazonaws.com/dbuggen2/marke.png"})
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
		issueID := pathIntSeparator(c.Param("issue"))
		articleIndex := pathIntSeparator(c.Param("article"))

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

func authortext(AuthorText sql.NullString, authors []database.Author) string {
	if AuthorText.Valid {
		return AuthorText.String
	}

	var sb strings.Builder
	sb.WriteString("Skriven av ")

	sb.WriteString(authorsName(authors[0]))
	if len(authors) == 1 {
		return sb.String()
	}

	for i := 1; i < len(authors)-1; i++ {
		sb.WriteString(fmt.Sprintf(", %v", authorsName(authors[i])))
	}

	sb.WriteString(fmt.Sprintf(" och %v", authorsName(authors[len(authors)-1])))
	return sb.String()
}

func authorsName(a database.Author) string {
	if a.PreferedName.Valid {
		return a.PreferedName.String
	}

	resp, err := http.Get(fmt.Sprintf("https://hodis.datasektionen.se/uid/%v", a.KthID))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected http GET status from hodis: %v", resp.StatusCode)
	}

	contents, err := io.ReadAll(io.Reader(resp.Body))
	if err != nil {
		log.Fatal(err)
	}
	m := make(map[string]interface{})
	errg := json.Unmarshal(contents, &m)
	if errg != nil {
		log.Fatal(err)
	}

	displayName, ok := m["displayName"].(string)
	if !ok {
		log.Fatal("Failed to convert displayName to string")
	}
	return displayName
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
