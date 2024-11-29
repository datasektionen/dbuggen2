package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func saveIssue(db *sqlx.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		order, err := determineOrder(c.PostForm("order-input"))
		if err != nil {
			panic("TODO")
		}
		log.Print(order)
	}
}
