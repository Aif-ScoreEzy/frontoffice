package loanrecordchecker

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/core/log/transaction"
	"front-office/pkg/core/product"
	"front-office/pkg/middleware"
	"front-office/pkg/procat/log"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(apiGroup fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repo := NewRepository(cfg, client, nil)
	productRepo := product.NewRepository(cfg, client)
	logRepo := log.NewRepository(cfg, client)
	transactionRepo := transaction.NewRepository(cfg, client)

	logService := log.NewService(logRepo, transactionRepo)
	service := NewService(repo, productRepo, logRepo, transactionRepo, logService)

	controller := NewController(service)

	loanRecordCheckerGroup := apiGroup.Group("loan-record-checker")
	loanRecordCheckerGroup.Post("/single-request", middleware.Auth(), middleware.IsRequestValid(loanRecordCheckerRequest{}), middleware.GetJWTPayloadFromCookie(), controller.LoanRecordChecker)
}
