package log

import (
	"front-office/app/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(logAPI fiber.Router,db *gorm.DB, cfg *config.Config) {
	repository := NewRepository(db, cfg)
	service := NewService(repository, cfg)
	controller := NewController(service)

	logAPI.Get("/", controller.GetTransactionLogs)
	logAPI.Get("/by-date", controller.GetTransactionLogsByDate)
	logAPI.Get("/by-range-date", controller.GetTransactionLogsByRangeDate)
	logAPI.Get("/by-month", controller.GetTransactionLogsByMonth)
}
