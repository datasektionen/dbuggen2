package client

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Home page
func Home(text string) func(c *gin.Context) {
	type Article struct {
		Name       string
		Date       time.Time
		Coverimage string
	}

	articles := []Article{{"loldbuggen", time.Now(), ""}, {"hejdbuggen", time.Now().Add(30), "bild.png"}}

	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":   "Main wow",
			"content": "home.tmpl",
			// Specific for home.tmpl ↓
			"articles": articles,
		})
	}
}
