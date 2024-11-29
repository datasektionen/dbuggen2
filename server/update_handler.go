package server

import (
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func determineOrder(str string) []int {
	orderStrings := strings.Split(str, ",")
	var order []int
	for i := 0; i < len(orderStrings); i++ {
		index, err := strconv.Atoi(orderStrings[i])
		if err != nil {
			continue
		}
		order = append(order, index)
	}
	return order
}

func saveIssue(db *sqlx.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		order := determineOrder(c.PostForm("order-input"))
		log.Print(order)
	}
}
