package main

import (
	"front-office/config"
	"front-office/config/database"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	loc := time.FixedZone("Asia/Jakarta", 25200)
	time.Local = loc

	app := fiber.New()

	config.LoadEnv()

	database.ConnectPostgres()

	log.Fatal(app.Listen(":" + os.Getenv("APP_PORT")))
}
