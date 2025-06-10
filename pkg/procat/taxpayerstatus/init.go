package taxpayerstatus

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(apiGroup fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repository := NewRepository(cfg, client)
	service := NewService(repository)
	controller := NewController(service)

	taxComplianceGroup := apiGroup.Group("tax-payer-status")
	taxComplianceGroup.Post("/", middleware.Auth(), middleware.IsRequestValid(taxPayerStatusRequest{}), middleware.GetJWTPayloadFromCookie(), controller.TaxPayerStatus)
}
