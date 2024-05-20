package main

import (
	// "dbuggen/server"
	"dbuggen/config"
	"dbuggen/server/database"
)

func main() {
	// server.Start()
	conf := config.GetConfig()
	database.DatabaseDo(conf.DATABASE_URL)
}
