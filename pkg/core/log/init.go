package log

import (
	"front-office/app/config"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupInit(logAPI fiber.Router, db *gorm.DB, cfg *config.Config) {
	repository := NewRepository(db, cfg)
	service := NewService(repository, cfg)
	controller := NewController(service)

	// scoreezy
	logTransScoreezyAPI := logAPI.Group("scoreezy")
	logTransScoreezyAPI.Get("/", controller.GetTransactionLogs)
	logTransScoreezyAPI.Get("/by-date", controller.GetTransactionLogsByDate)
	logTransScoreezyAPI.Get("/by-range-date", controller.GetTransactionLogsByRangeDate)
	logTransScoreezyAPI.Get("/by-month", controller.GetTransactionLogsByMonth)
}
