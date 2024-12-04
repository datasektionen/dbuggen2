package server

import (
	"database/sql"
	"dbuggen/server/database"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func determineOrder(str string) ([]int, error) {
	orderStrings := strings.Split(str, ",")
	var order []int
	for i := 0; i < len(orderStrings); i++ {
		index, err := strconv.Atoi(orderStrings[i])
		if err != nil {
			return nil, err
		}
		order = append(order, index)
	}
	return order, nil
}

func extractArticleData(c *gin.Context, articleID int, index int, issue int) (database.Article, error) {
	articleIDString := strconv.Itoa(articleID)
	title := c.PostForm(articleIDString + "_Title")
	// authors := c.PostForm(articleIDString + "_Authors")
	authortext := sql.NullString{
		String: c.PostForm(articleIDString + "_Authortext"),
		Valid:  true,
	}
	// N0lleSafe := c.PostForm(articleIDString + "_NØllesafe")
	content := c.PostForm(articleIDString + "_Content")

	article := database.Article{
		ID:         articleID,
		Title:      title,
		Issue:      issue,
		AuthorText: authortext,
		IssueIndex: index,
		Content:    content,
		LastEdited: time.Now(),
		N0lleSafe:  true,
	}
	return article, nil
}

func checkArticleInIssue(articleIDs []int, currentArticles []database.Article) ([]int, error) {
	editedArticles := make([]bool, len(currentArticles))
	for _, aID := range articleIDs {
		if aID < 0 {
			continue
		}

		inIssue := false
		for i, article := range currentArticles {
			if aID == article.ID {
				inIssue = true
				editedArticles[i] = true
				break
			}
		}

		if !inIssue {
			return nil, fmt.Errorf("articleID %v is not in current issue and may not be edited", aID)
		}
	}

	var articlesRemove []int
	for i, b := range editedArticles {
		if !b {
			articlesRemove = append(articlesRemove, currentArticles[i].ID)
		}
	}

	return articlesRemove, nil
}
