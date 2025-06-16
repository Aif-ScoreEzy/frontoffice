package loanrecordchecker

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

	loanRecordCheckerGroup := apiGroup.Group("loan-record-checker")
	loanRecordCheckerGroup.Post("/single-request", middleware.Auth(), middleware.IsRequestValid(LoanRecordCheckerRequest{}), middleware.GetJWTPayloadFromCookie(), controller.LoanRecordChecker)
}
