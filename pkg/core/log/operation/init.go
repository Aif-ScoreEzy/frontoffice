package operation

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(logAPI fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repository := NewRepository(cfg, client)
	service := NewService(repository)
	controller := NewController(service)

	logOperationAPI := logAPI.Group("operation")
	logOperationAPI.Get("/", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetList)
	logOperationAPI.Get("/range", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.GetListByRange)
}
