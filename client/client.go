package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Generate the html for the home page
func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":   "Main wow",
		"content": "home.tmpl",
		// Specific for home.tmpl ↓
		"text": "jag älskar faktiskt kth",
	})
}
