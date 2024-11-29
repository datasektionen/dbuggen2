package server

import (
	"database/sql"
	"dbuggen/server/database"
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

func extractArticleData(c *gin.Context, articleID int) (database.Article, error) {
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
		Issue:      -1,
		AuthorText: authortext,
		IssueIndex: -1,
		Content:    content,
		LastEdited: time.Now(),
		N0lleSafe:  true,
	}
	return article, nil
}
