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
	repo := NewRepository(cfg, client)
	productRepo := product.NewRepository(cfg, client)
	logRepo := log.NewRepository(cfg, client)
	transRepo := transaction.NewRepository(cfg, client)

	service := NewService(repo, productRepo, logRepo, transRepo)
	productService := product.NewService(productRepo)
	logService := log.NewService(logRepo)
	transService := transaction.NewService(transRepo)

	controller := NewController(service, productService, logService, transService)

	loanRecordCheckerGroup := apiGroup.Group("loan-record-checker")
	loanRecordCheckerGroup.Post("/single-request", middleware.Auth(), middleware.IsRequestValid(LoanRecordCheckerRequest{}), middleware.GetJWTPayloadFromCookie(), controller.LoanRecordChecker)
}
