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

	apiGroup.Post("/7d-multiple-loan/single-request", middleware.Auth(), middleware.IsRequestValid(multipleLoanRequest{}), middleware.GetJWTPayloadFromCookie(), controller.MultipleLoan7Days)
	apiGroup.Post("/30d-multiple-loan/single-request", middleware.Auth(), middleware.IsRequestValid(multipleLoanRequest{}), middleware.GetJWTPayloadFromCookie(), controller.MultipleLoan30Days)
	apiGroup.Post("/90d-multiple-loan/single-request", middleware.Auth(), middleware.IsRequestValid(multipleLoanRequest{}), middleware.GetJWTPayloadFromCookie(), controller.MultipleLoan90Days)
}
