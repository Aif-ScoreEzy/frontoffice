package loanrecordchecker

import (
	"front-office/app/config"
	"front-office/internal/httpclient"
	"front-office/pkg/core/log/transaction"
	"front-office/pkg/core/product"
	"front-office/pkg/middleware"
	"front-office/pkg/procat/job"

	"github.com/gofiber/fiber/v2"
)

func SetupInit(apiGroup fiber.Router, cfg *config.Config, client httpclient.HTTPClient) {
	repo := NewRepository(cfg, client, nil)
	productRepo := product.NewRepository(cfg, client)
	jobRepo := job.NewRepository(cfg, client, nil)
	transactionRepo := transaction.NewRepository(cfg, client, nil)

	jobService := job.NewService(jobRepo, transactionRepo)
	service := NewService(repo, productRepo, jobRepo, transactionRepo, jobService)

	controller := NewController(service)

	loanRecordCheckerGroup := apiGroup.Group("loan-record-checker")
	loanRecordCheckerGroup.Post("/single-request", middleware.Auth(), middleware.IsRequestValid(loanRecordCheckerRequest{}), middleware.GetJWTPayloadFromCookie(), controller.SingleSearch)
	loanRecordCheckerGroup.Post("/bulk-request", middleware.Auth(), middleware.GetJWTPayloadFromCookie(), controller.BulkSearch)
}
