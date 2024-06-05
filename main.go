package main

import (
	"dbuggen/config"
	"dbuggen/server"
	"dbuggen/server/database"
)

func main() {
	conf := config.GetConfig()
	db := database.Start(conf.DATABASE_URL)
	server.Start(db, conf)
}
