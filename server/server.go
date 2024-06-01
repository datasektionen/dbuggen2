package server

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"dbuggen/client"
	"dbuggen/server/database"
)

// Start starts the server and initializes the routes and templates.
func Start(db *sqlx.DB) {
	r := gin.Default()
	tmpl := template.Must(template.ParseGlob("client/html/*.tmpl"))
	r.SetHTMLTemplate(tmpl)

	r.Static("css", "client/css")
	r.Static("assets", "assets")

	r.GET("/", client.Home(database.GetIssues(db)))
	r.GET("issue/:issue/*article", client.Article(db))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.Run()
}

// func article(db *sqlx.DB) func(c *gin.Context) {
// 	return client.Article(db)
// }
