package client

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"dbuggen/server/database"
)

//go:embed html/*.html
var HTMLTemplates embed.FS

//go:embed public/*
var PublicFiles embed.FS

// Home page
func Home(db *sqlx.DB, ds *DarkmodeStatus) func(c *gin.Context) {
	return func(c *gin.Context) {
		issuesRaw, err := database.GetHomeIssues(db, Darkmode(ds))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
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
			c.Redirect(http.StatusBadRequest, "/")
			return
		}

		darkmode := Darkmode(ds)

		issue, err := database.GetIssue(db, issueID, darkmode)
		if err != nil {
			c.Redirect(http.StatusInternalServerError, "/")
			return
		}

		articles, err := database.GetArticles(db, issueID, darkmode)
		if err != nil {
			c.Redirect(http.StatusInternalServerError, "/")
			return
		}

		databaseAuthors, err := database.GetAuthorsForIssue(db, issueID)
		if err != nil {
			c.Redirect(http.StatusInternalServerError, "/")
			return
		}

		var issueArticles []issueArticle
		for _, article := range articles {
			var authors string
			if len(databaseAuthors) <= article.IssueIndex {
				var a []database.Author
				authors = authortext(article.AuthorText, a)
			} else {
				authors = authortext(article.AuthorText, databaseAuthors[article.IssueIndex])
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
		authors, err := database.GetAuthorsForArticle(db, article.ID)
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
