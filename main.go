package main

import (
	"fmt"
	"front-office/config"
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

	fmt.Println(os.Getenv("APP_ENV"))

	log.Fatal(app.Listen(":" + os.Getenv("APP_PORT")))
}
