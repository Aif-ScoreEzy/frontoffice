package loanrecordchecker

import (
	"front-office/app/config"
	"front-office/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(apiGroup fiber.Router, cfg *config.Config) {
	repository := NewRepository(cfg)
	service := NewService(cfg, repository)
	controller := NewController(service)

	loanRecordCheckerGroup := apiGroup.Group("loan-record-checker")
	loanRecordCheckerGroup.Post("/", middleware.Auth(), middleware.IsRequestValid(LoanRecordCheckerRequest{}), middleware.GetJWTPayloadFromCookie(), controller.LoanRecordChecker)
}
