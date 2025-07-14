package transaction

import (
	"front-office/app/config"
	"front-office/internal/httpclient"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(logAPI fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repository := NewRepository(cfg, client, nil)
	service := NewService(repository)
	controller := NewController(service)

	logTransScoreezyAPI := logAPI.Group("scoreezy")
	logTransScoreezyAPI.Get("/", controller.GetLogScoreezy)
	logTransScoreezyAPI.Get("/by-date", controller.GetLogScoreezyByDate)
	logTransScoreezyAPI.Get("/by-range-date", controller.GetLogScoreezyByDateRange)
	logTransScoreezyAPI.Get("/by-month", controller.GetLogScoreezyByMonth)

	// logTransProcatGroup := logAPI.Group("product_catalog")
}
