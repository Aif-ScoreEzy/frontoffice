package multipleloan

import (
	"front-office/app/config"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(apiGroup fiber.Router, cfg *config.Config) {
	repository := NewRepository(cfg)
	service := NewService(cfg, repository)
	controller := NewController(service)

	multipleLoanGroup := apiGroup.Group("multiple-loan")
	multipleLoanGroup.Post("/7-days", middleware.Auth(), middleware.IsRequestValid(MultipleLoanRequest{}), middleware.GetJWTPayloadFromCookie(), controller.MultipleLoan7Days)
	multipleLoanGroup.Post("/30-days", middleware.Auth(), middleware.IsRequestValid(MultipleLoanRequest{}), middleware.GetJWTPayloadFromCookie(), controller.MultipleLoan30Days)
}
