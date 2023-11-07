package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	if env == "local" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalln("Error loading .env file")
		}
	}

}
