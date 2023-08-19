package main

import (
	"front-office/config"
	database "front-office/config/database"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	loc := time.FixedZone("Asia/Jakarta", 25200)
	time.Local = loc

	app := fiber.New()
	app.Use(recover.New())

	config.LoadEnv()

	database.ConnectPostgres()
	config.Migrate()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	config.SetupRoutes(app)
	config.SetupRoutes(app)

	log.Fatal(app.Listen(":" + os.Getenv("APP_PORT")))
}
