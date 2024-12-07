package server

import (
	"dbuggen/client"
	"dbuggen/server/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func saveIssue(db *sqlx.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		order, err := determineOrder(c.PostForm("order-input"))
		if err != nil {
			c.JSON(http.StatusBadRequest, "")
			return
		}

		issue, err := client.PathIntSeparator(c.Param("issue"))
		if err != nil {
			c.JSON(http.StatusBadRequest, "")
			return
		}

		currentArticles, err := database.GetArticles(db, issue, false)
		if err != nil {
			c.JSON(http.StatusBadRequest, "")
			return
		}

		articlesRemove, err := checkArticleInIssue(order, currentArticles)
		if err != nil {
			c.JSON(http.StatusBadRequest, "")
			return
		}

		var articlesModify []database.Article
		var articlesNew []database.Article
		for i, a := range order {
			article, err := extractArticleData(c, a, i, issue)
			if err != nil {
				c.JSON(http.StatusBadRequest, "")
				return
			}
			if a >= 0 {
				articlesModify = append(articlesModify, article)
			} else {
				articlesNew = append(articlesNew, article)
			}
		}

		err = database.UpdateIssue(db, articlesModify, articlesNew, articlesRemove)
		if err != nil {
			c.JSON(http.StatusBadRequest, "")
			return
		}
	}
}

func deleteIssue(db *sqlx.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		issue, err := client.PathIntSeparator(c.Param("issue"))
		if err != nil {
			c.JSON(http.StatusBadRequest, "")
			return
		}
		err = database.DeleteIssue(db, issue)
		if err != nil {
			c.JSON(http.StatusBadRequest, "")
			return
		}
	}
}
