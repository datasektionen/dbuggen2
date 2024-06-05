package client

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
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

// coverpage generates an HTML template for the cover image.
// If the coverpage is valid, it returns an HTML string with an image tag.
// If the coverpage is not valid, it returns an empty HTML string.
func coverpage(coverpage sql.NullString) template.HTML {
	if coverpage.Valid {
		return template.HTML(fmt.Sprintf(`<img src="%v" style="max-width: 40vw;">`, template.HTMLEscapeString(coverpage.String)))
	}
	return template.HTML("")
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
func pathIntSeparator(paramRaw string) (int, error) {
	_, paramLessRaw, found := strings.Cut(paramRaw, "/")

	if found {
		param, err := strconv.Atoi(paramLessRaw)
		if err != nil {
			return 0, err
		}
		return param, nil
	}

	param, err := strconv.Atoi(paramRaw)
	if err != nil {
		return 0, err
	}
	return param, nil
}
