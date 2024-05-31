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
	r.Static("/css", "client/css")
	r.Static("assets", "assets")

	r.GET("/", client.Home(database.GetIssues(db)))

	r.SetHTMLTemplate(tmpl)
	r.Run()
}
