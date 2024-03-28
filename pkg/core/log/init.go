package log

import (
	"front-office/app/config"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(logAPI fiber.Router, cfg *config.Config) {
	service := NewService(cfg)
	controller := NewController(service)

	logAPI.Get("/by-date", controller.GetTransactionLogsByDate)
	logAPI.Get("/by-range-date", controller.GetTransactionLogsByRangeDate)
	logAPI.Get("/by-month", controller.GetTransactionLogsByMonth)
	logAPI.Get("/by-name", controller.GetTransactionLogsByName)
}
