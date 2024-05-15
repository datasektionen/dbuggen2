package server

import (
	"html/template"

	"github.com/gin-gonic/gin"

	"dbuggen/client"
)

// Start starts the server and initializes the routes and templates.
func Start() {
	r := gin.Default()

	tmpl := template.Must(template.ParseGlob("client/html/*.tmpl"))
	r.Static("/css", "client/css")
	r.Static("assets", "assets")

	r.GET("/", client.Home("lol vad sjuktr"))

	r.SetHTMLTemplate(tmpl)
	r.Run()
}
