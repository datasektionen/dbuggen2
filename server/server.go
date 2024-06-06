package server

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"dbuggen/client"
	"dbuggen/config"
)

// Start starts the server and initializes the routes and templates.
func Start(db *sqlx.DB, conf *config.Config) {
	r := gin.Default()
	tmpl := template.Must(template.ParseGlob("client/html/*.tmpl"))
	r.SetHTMLTemplate(tmpl)

	r.Static("css", "client/css")
	r.Static("assets", "assets")

	ds := initDarkmode(conf.DARKMODE_URL)
	r.GET("/", client.Home(db, &ds, conf.DARKMODE_URL))
	r.GET("issue/:issue/:article", client.Article(db))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.Run()
}

func initDarkmode(url string) client.DarkmodeStatus {
	ds := client.DarkmodeStatus{
		Darkmode: true,
		LastPoll: time.Date(1983, time.October, 7, 17, 0, 0, 0, time.Local),
	}

	client.Darkmode(&ds, url)
	return ds
}
