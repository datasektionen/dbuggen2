package server

import (
	"html/template"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"dbuggen/client"
	"dbuggen/config"
)

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

// Start starts the server and initializes the routes and templates.
func Start(db *sqlx.DB, conf *config.Config) {
	r := gin.Default()
	r.SetHTMLTemplate(template.Must(template.ParseFS(client.HTMLTemplates, "**/*.html")))

	r.StaticFS("public", http.FS(must(fs.Sub(client.PublicFiles, "public"))))

	var ds client.DarkmodeStatus
	initDarkmode(&ds, conf.DARKMODE_URL)

	r.GET("/", client.Home(db, &ds))
	r.GET("issue/:issue", client.Issue(db, &ds))
	r.GET("issue/:issue/edit", client.EditIssue(db))
	r.GET("issue/:issue/:article", client.Article(db, &ds))
	r.GET("redaqtionen", client.Redaqtionen(db, conf.DFUNKT_URL))

	r.POST("issue/:issue/save", saveIssue(db))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	r.Run()
}

func initDarkmode(ds *client.DarkmodeStatus, url string) {
	*ds = client.DarkmodeStatus{
		Darkmode: true,
		LastPoll: time.Date(1983, time.October, 7, 17, 0, 0, 0, time.Local),
		Url:      url,
		Mutex:    sync.RWMutex{},
	}

	client.Darkmode(ds)
}
