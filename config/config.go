package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DATABASE_URL string
}

func GetConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println(err)
	}

	conf := Config{
		DATABASE_URL: os.Getenv("DATABASE_URL"),
	}

	return &conf
}
