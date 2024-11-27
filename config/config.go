package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DATABASE_URL string
	DFUNKT_URL   string
	DARKMODE_URL string
}

func GetConfig() *Config {
	_ = godotenv.Load()

	conf := Config{
		DATABASE_URL: os.Getenv("DATABASE_URL"),
		DFUNKT_URL:   os.Getenv("DFUNKT_URL"),
		DARKMODE_URL: os.Getenv("DARKMODE_URL"),
	}

	return &conf
}
