package multipleloan

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(apiGroup fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repository := NewRepository(cfg, client)
	service := NewService(cfg, repository)
	controller := NewController(service)

	multipleLoanGroup := apiGroup.Group("multiple-loan")
	multipleLoanGroup.Post("/7-days", middleware.Auth(), middleware.IsRequestValid(MultipleLoanRequest{}), middleware.GetJWTPayloadFromCookie(), controller.MultipleLoan7Days)
	multipleLoanGroup.Post("/30-days", middleware.Auth(), middleware.IsRequestValid(MultipleLoanRequest{}), middleware.GetJWTPayloadFromCookie(), controller.MultipleLoan30Days)
	multipleLoanGroup.Post("/90-days", middleware.Auth(), middleware.IsRequestValid(MultipleLoanRequest{}), middleware.GetJWTPayloadFromCookie(), controller.MultipleLoan90Days)
}
