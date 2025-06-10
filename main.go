package main

import (
	"front-office/app/config"

	"front-office/app/server"
	"time"
)

func main() {
	loc := time.FixedZone("Asia/Jakarta", 25200)
	time.Local = loc

	cfg := config.GetConfig()

	// migrate.PostgreDB(db)

	server.NewServer(&cfg).Start()
}
